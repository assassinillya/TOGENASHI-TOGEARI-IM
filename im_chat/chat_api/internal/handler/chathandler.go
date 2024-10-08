package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"gorm.io/gorm"
	"im_server/common/models/ctype"
	"im_server/common/response"
	"im_server/common/service/redis_service"
	"im_server/im_chat/chat_api/internal/svc"
	"im_server/im_chat/chat_api/internal/types"
	"im_server/im_chat/chat_models"
	"im_server/im_file/file_rpc/types/file_rpc"
	"im_server/im_user/user_models"
	"im_server/im_user/user_rpc/types/user_rpc"
	"net/http"
	"strings"
	"time"
)

type UserWsInfo struct {
	UserInfo    user_models.UserModel      //用户信息
	WsClientMap map[string]*websocket.Conn // 这个用户管理的所有ws客户端
	CurrenConn  *websocket.Conn            // 当前的连接对象
}

var UserOnlineWsMap = map[uint]*UserWsInfo{}

func chatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatRequest
		if err := httpx.ParseHeaders(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// todo 鉴权 true表示放行，false表示拦截
				return true
			},
		}

		conn, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		addr := conn.RemoteAddr().String()
		defer func() {
			conn.Close()

			userWsInfo, ok := UserOnlineWsMap[req.UserID]
			if ok {
				// 删除的是推出的那个ws信息
				delete(userWsInfo.WsClientMap, addr)
			}

			if userWsInfo != nil && len(userWsInfo.WsClientMap) == 0 {
				// 用户全退出了
				delete(UserOnlineWsMap, req.UserID)
				svcCtx.Redis.HDel("online", fmt.Sprintf("%d", req.UserID))
			}

		}()
		//调用户服务，获取当前用户信息
		res, err := svcCtx.UserRpc.UserInfo(context.Background(), &user_rpc.UserInfoRequest{
			UserId: uint32(req.UserID),
		})
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		var userInfo user_models.UserModel
		json.Unmarshal(res.Data, &userInfo)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		userWsInfo, ok := UserOnlineWsMap[req.UserID]
		if !ok {
			// 代表这个用户第一次来
			userWsInfo = &UserWsInfo{
				UserInfo: userInfo,
				WsClientMap: map[string]*websocket.Conn{
					addr: conn,
				},
				CurrenConn: conn, // 当前的连接对象
			}
			UserOnlineWsMap[req.UserID] = userWsInfo
		}
		_, ok1 := userWsInfo.WsClientMap[addr]
		if !ok1 {
			// 代表这个用户二开及以上
			UserOnlineWsMap[req.UserID].WsClientMap[addr] = conn
			// 把当前连接对象更换
			UserOnlineWsMap[req.UserID].CurrenConn = conn
		}

		// 把在线的用户存入redis
		svcCtx.Redis.HSet("online", fmt.Sprintf("%d", req.UserID), req.UserID)

		// 遍历在线的用户, 如果与当前用户是好友, 就给他发好友在线

		// 先把所有在线的用户id取出来, 以及待确认的用户id, 然后传到用户rpc服务中
		// [1,2,3]  3
		// 在rpc服务, 去判断哪些用户是好友关系

		// 如果好友开启了好友上线提醒
		// 查一下自己的好友是不是上线了
		friendsRes, err := svcCtx.UserRpc.FriendList(context.Background(), &user_rpc.FriendListRequest{
			User: uint32(req.UserID),
		})
		// 3 [3,4,5]
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		logx.Infof("用户上线: %s 用户id: %d", userInfo.Nickname, req.UserID)

		for _, info := range friendsRes.FriendList {
			friend, ok := UserOnlineWsMap[uint(info.UserId)]
			if ok {
				text := fmt.Sprintf("好友%s上线了", UserOnlineWsMap[req.UserID].UserInfo.Nickname)
				// 判断用户是否开了好友上线提示功能
				if friend.UserInfo.UserConfModel.FriendOnline {
					// 好友上线了
					//friend.Conn.WriteMessage(websocket.TextMessage, []byte(text))
					sendWsMapMsg(friend.WsClientMap, []byte(text))
				}

			}
		}
		// 查一下自己的好友列表, 返回用户id列表, 看看UserWsMap中是否存在, 如果存在就给自己发一个好友上线的消息

		logx.Info(UserOnlineWsMap)

		defer conn.Close()
		for {
			// 消息类型，消息，错误
			_, p, err1 := conn.ReadMessage()
			if err1 != nil {
				// 用户断开聊天
				fmt.Println(err1)
				break
			}
			var request ChatRequest
			err2 := json.Unmarshal(p, &request)
			if err2 != nil {
				// 格式不正确
				logx.Error(err2)

				SendTipErrMsg(conn, "参数解析失败")
				continue
			}
			if request.RevUserID != req.UserID {
				// 判断你聊天的这个人是不是你的好友
				isFriendRes, err3 := svcCtx.UserRpc.IsFriend(context.Background(), &user_rpc.IsFriendRequest{
					User1: uint32(req.UserID),
					User2: uint32(request.RevUserID),
				})
				if err3 != nil {
					logx.Error(err3)
					SendTipErrMsg(conn, "用户服务错误, 请重试")
					return
				}

				if !isFriendRes.IsFriend {
					SendTipErrMsg(conn, "你们还不是好友")
					return
				}
			}

			// 判断type 1-12
			if !(request.Msg.Type >= 1 && request.Msg.Type <= 12) {
				SendTipErrMsg(conn, "消息类型错误")
				continue
			}

			// 判断是否是文件类型
			switch request.Msg.Type {
			case ctype.TextMsgType:
				if request.Msg.TextMsg == nil {
					SendTipErrMsg(conn, "请输入消息内容")
					continue
				}

				if request.Msg.TextMsg.Content == "" {
					SendTipErrMsg(conn, "请输入消息内容")
					continue
				}
			case ctype.FileMsgType:

				if request.Msg.FileMsg == nil {
					SendTipErrMsg(conn, "请上传文件")
					return
				}

				nameList := strings.Split(request.Msg.FileMsg.Src, "/")
				if len(nameList) == 0 {
					SendTipErrMsg(conn, "请上传文件")
					continue
				}

				fileID := nameList[len(nameList)-1]
				fileResponse, err4 := svcCtx.FileRpc.FileInfo(context.Background(), &file_rpc.FileInfoRequest{
					FileId: fileID,
				})
				if err4 != nil {
					logx.Error(err4)
					SendTipErrMsg(conn, err4.Error())
					continue
				}
				request.Msg.FileMsg.Title = fileResponse.FileName
				request.Msg.FileMsg.Size = fileResponse.FileSize
				request.Msg.FileMsg.Type = fileResponse.FileType
			case ctype.WithdrawMsgType:
				// 撤回消息的消息id是必填的
				if request.Msg.WithdrawMsg == nil {
					SendTipErrMsg(conn, "撤回消息id必填")
					return
				}

				if request.Msg.WithdrawMsg.MsgID == 0 {
					SendTipErrMsg(conn, "撤回消息id必填")
					continue
				}

				// 自己只能撤回自己的消息
				// 查找这个消息是谁发的
				var msgModel chat_models.ChatModel
				err = svcCtx.DB.Take(&msgModel, request.Msg.WithdrawMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "消息不存在")
					continue
				}

				//不能撤回已撤回的消息
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "撤回消息不能再撤回了")
					continue
				}

				// 查找这个信息是不是自己发的
				if msgModel.SendUserID != req.UserID {
					SendTipErrMsg(conn, "只能撤回自己的消息")
					continue
				}

				// 判断消息的时候, 小于2分钟的消息不能撤回
				if time.Now().Sub(msgModel.CreatedAt) >= time.Minute*2 {
					SendTipErrMsg(conn, "只能撤回2分钟以内的消息")
					continue
				}

				// 撤回逻辑
				// 收到撤回请求之后, 服务端把原消息的类型改为撤回消息类型, 并且记录原消息
				// 然后通知前端的收发方, 重新拉取聊天记录

				var content = "撤回了一条消息"
				if userInfo.UserConfModel.RecallMessage != nil {
					content = *userInfo.UserConfModel.RecallMessage
				}
				content = "你" + content

				// 前端可以判断, 这个消息如果不是isMe, 就可以把你替换成对方的昵称
				// 这里涉及一个循环引用的问题, originMsg里也包括了撤回信息导致无限循环
				originMsg := msgModel.Msg
				originMsg.WithdrawMsg = nil // 这里可能会循环引用, 所以拷贝了这个值并且把撤回消息置空了
				svcCtx.DB.Model(&msgModel).Updates(chat_models.ChatModel{
					MsgPreView: "[撤回消息] - " + content,
					MsgType:    ctype.WithdrawMsgType,
					Msg: ctype.Msg{
						Type: ctype.WithdrawMsgType,
						WithdrawMsg: &ctype.WithdrawMsg{
							Content:   content,
							MsgID:     request.Msg.WithdrawMsg.MsgID,
							OriginMsg: &originMsg,
						},
					},
				})

				// 把原消息置空
			case ctype.ReplyMsgType:
				//回复消息
				//先校验
				if request.Msg.ReplyMsg == nil || request.Msg.ReplyMsg.MsgID == 0 {
					SendTipErrMsg(conn, "回复消息必填")
					continue
				}
				//找这个原消息
				var msgModel chat_models.ChatModel
				err = svcCtx.DB.Take(&msgModel, request.Msg.ReplyMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "消息不存在")
					continue
				}

				// 不能回复撤回消息
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "该消息已撤回")
					continue
				}

				// 回复的消息必须是你自己或者和你聊天的这个人发出来的

				// 原消息必须是 当前你要和对方聊的 原消息就会有一个发送人id和接收人id, 我们的聊天也会有一个发送人id和接收人id
				// 因为回复消息可以回复自己的, 也可以回复别人的
				// 这里注意打开的会话必须是与别人的对话才能回复, 如果和自己对话回复别人的对话是不可以的
				if !((msgModel.SendUserID == req.UserID && msgModel.RevUserID == request.RevUserID) ||
					(msgModel.SendUserID == request.RevUserID && msgModel.RevUserID == req.UserID)) {
					SendTipErrMsg(conn, "只能回复自己或者对方的信息")
					continue
				}

				userBaseInfo, err3 := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, msgModel.SendUserID)
				if err3 != nil {
					logx.Error(err3)
					SendTipErrMsg(conn, err3.Error())
					continue
				}

				request.Msg.ReplyMsg.Msg = &msgModel.Msg
				request.Msg.ReplyMsg.UserID = msgModel.SendUserID
				request.Msg.ReplyMsg.UserNickName = userBaseInfo.NickName
				request.Msg.ReplyMsg.OriginMsgDate = msgModel.CreatedAt
			case ctype.QuoteMsgType:
				// 引用消息
				// 先校验
				if request.Msg.QuoteMsg == nil || request.Msg.QuoteMsg.MsgID == 0 {
					SendTipErrMsg(conn, "引用消息必填")
					continue
				}
				//找这个原消息
				var msgModel chat_models.ChatModel
				err = svcCtx.DB.Take(&msgModel, request.Msg.QuoteMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "消息不存在")
					continue
				}

				// 不能回复撤回消息
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "该消息已撤回")
					continue
				}

				// 回复的消息必须是你自己或者和你聊天的这个人发出来的

				// 原消息必须是 当前你要和对方聊的 原消息就会有一个发送人id和接收人id, 我们的聊天也会有一个发送人id和接收人id
				// 因为回复消息可以回复自己的, 也可以回复别人的
				// 这里注意打开的会话必须是与别人的对话才能回复, 如果和自己对话回复别人的对话是不可以的
				if !((msgModel.SendUserID == req.UserID && msgModel.RevUserID == request.RevUserID) ||
					(msgModel.SendUserID == request.RevUserID && msgModel.RevUserID == req.UserID)) {
					SendTipErrMsg(conn, "只能回复自己或者对方的信息")
					continue
				}

				userBaseInfo, err3 := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, msgModel.SendUserID)
				if err3 != nil {
					logx.Error(err3)
					SendTipErrMsg(conn, err3.Error())
					continue
				}

				request.Msg.QuoteMsg.Msg = &msgModel.Msg
				request.Msg.QuoteMsg.UserID = msgModel.SendUserID
				request.Msg.QuoteMsg.UserNickName = userBaseInfo.NickName
				request.Msg.QuoteMsg.OriginMsgDate = msgModel.CreatedAt
			}

			// 先入库
			msgID := InsertMsgByChat(svcCtx.DB, req.UserID, request.RevUserID, request.Msg)
			// 判断目标用户在不在线 给发送双方都要发消息
			SendMsgByUser(svcCtx, req.UserID, request.RevUserID, request.Msg, msgID)
		}

	}
}

type ChatRequest struct {
	RevUserID uint      `json:"revUserID"` // 给谁发
	Msg       ctype.Msg `json:"msg"`
}

func sendWsMapMsg(wsMap map[string]*websocket.Conn, byteData []byte) {
	for _, conn := range wsMap {
		conn.WriteMessage(websocket.TextMessage, byteData)
	}
}

type ChatResponse struct {
	ID        uint           `json:"id"`
	IsMe      bool           `json:"isMe"`
	RevUser   ctype.UserInfo `json:"revUser"`
	SendUser  ctype.UserInfo `json:"sendUser"`
	Msg       ctype.Msg      `json:"msg"`
	CreatedAt time.Time      `json:"created_at"`
}

// InsertMsgByChat 消息入库
func InsertMsgByChat(db *gorm.DB, sendUserID uint, revUserID uint, msg ctype.Msg) (msgID uint) {

	switch msg.Type {
	case ctype.WithdrawMsgType:
		logx.Info("撤回消息是不自己入库的")
		return
	}

	chatModel := chat_models.ChatModel{
		SendUserID: sendUserID,
		RevUserID:  revUserID,
		MsgType:    msg.Type,
		Msg:        msg,
	}
	chatModel.MsgPreView = chatModel.MsgPreviewMethod()
	err := db.Create(&chatModel).Error
	if err != nil {
		logx.Error(err)
		sendUser, ok := UserOnlineWsMap[sendUserID]
		if !ok {
			return
		}
		// todo 只会给最新的人, 发送错误消息, 但不一定是最新的人报的错
		SendTipErrMsg(sendUser.CurrenConn, "消息保存失败")

	}
	return chatModel.ID
}

// SendMsgByUser 发消息 给谁发 谁发的
func SendMsgByUser(svcCtx *svc.ServiceContext, sendUserID uint, revUserID uint, msg ctype.Msg, msgID uint) {

	revUser, ok1 := UserOnlineWsMap[revUserID]
	sendUser, ok2 := UserOnlineWsMap[sendUserID]

	resp := ChatResponse{
		ID:        msgID,
		Msg:       msg,
		CreatedAt: time.Now(),
	}

	// 自己与自己发消息
	if ok1 && ok2 && sendUserID == revUserID {
		//百分百是自己与自己发消息了
		resp.SendUser = ctype.UserInfo{
			ID:       sendUserID,
			NickName: sendUser.UserInfo.Nickname,
			Avatar:   sendUser.UserInfo.Avatar,
		}
		resp.RevUser = ctype.UserInfo{
			ID:       revUserID,
			NickName: revUser.UserInfo.Nickname,
			Avatar:   revUser.UserInfo.Avatar,
		}
		byteData, _ := json.Marshal(resp)
		//revUser.Conn.WriteMessage(websocket.TextMessage, byteData)
		sendWsMapMsg(revUser.WsClientMap, byteData)
		return
	}

	// 在线的情况下, 我是可以拿到对方的在线信息的
	// 对方不在线的情况下, 我只能通过调用用户服务rpc方法来获取

	// 无论如何都要给发送者回传消息

	if !ok1 {
		userBaseInfo, err := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, revUserID)

		if err != nil {
			logx.Error(err)
			return
		}

		resp.RevUser = ctype.UserInfo{
			ID:       revUserID,
			NickName: userBaseInfo.NickName,
			Avatar:   userBaseInfo.Avatar,
		}
	} else {
		resp.RevUser = ctype.UserInfo{
			ID:       revUserID,
			NickName: revUser.UserInfo.Nickname,
			Avatar:   revUser.UserInfo.Avatar,
		}
	}

	// 发送者在线
	resp.SendUser = ctype.UserInfo{
		ID:       sendUserID,
		NickName: sendUser.UserInfo.Nickname,
		Avatar:   sendUser.UserInfo.Avatar,
	}
	resp.IsMe = true

	byteData, _ := json.Marshal(resp)
	//sendUser.Conn.WriteMessage(websocket.TextMessage, byteData)
	sendWsMapMsg(sendUser.WsClientMap, byteData)

	//if ok1 && ok2 && sendUserID == revUserID{ // 自己加的判断是否为自己
	//	resp.IsMe = true
	//	byteData, _ = json.Marshal(resp)
	//	revUser.Conn.WriteMessage(websocket.TextMessage, byteData)
	//} else
	if ok1 {
		// 接受者在线
		resp.IsMe = false
		byteData, _ = json.Marshal(resp)
		//revUser.Conn.WriteMessage(websocket.TextMessage, byteData)
		sendWsMapMsg(revUser.WsClientMap, byteData)

	}

}

// SendTipErrMsg 发送错误提示的消息
func SendTipErrMsg(conn *websocket.Conn, msg string) {
	resp := ChatResponse{
		Msg: ctype.Msg{
			Type: ctype.TipMsgType,
			TipMsg: &ctype.TipMsg{
				Status:  "error",
				Content: msg,
			},
		},
		CreatedAt: time.Now(),
	}
	byteData, _ := json.Marshal(resp)
	conn.WriteMessage(websocket.TextMessage, byteData)
}
