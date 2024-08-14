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
	UserInfo user_models.UserModel //用户信息
	Conn     *websocket.Conn       // 用户的ws连接对象
}

var UserOnlineWsMap = map[uint]UserWsInfo{}

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

		defer func() {
			conn.Close()
			delete(UserOnlineWsMap, req.UserID)
			svcCtx.Redis.HDel("online", fmt.Sprintf("%d", req.UserID))
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

		var userWsInfo = UserWsInfo{
			UserInfo: userInfo,
			Conn:     conn,
		}

		UserOnlineWsMap[req.UserID] = userWsInfo
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
					friend.Conn.WriteMessage(websocket.TextMessage, []byte(text))
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
				if request.Msg.WithdrawMsg.MsgID == 0 {
					SendTipErrMsg(conn, "撤回消息id必填")
					return
				}

				// 自己只能撤回自己的消息
				// 查找这个消息是谁发的
				var msgModel chat_models.ChatModel
				err = svcCtx.DB.Take(&msgModel, request.Msg.WithdrawMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "消息不存在")
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

				var content = "xx 撤回了一条消息"
				if userInfo.UserConfModel.RecallMessage != nil {
					content = *userInfo.UserConfModel.RecallMessage
				}

				svcCtx.DB.Model(&msgModel).Updates(chat_models.ChatModel{
					Msg: ctype.Msg{
						Type: ctype.WithdrawMsgType,
						WithdrawMsg: &ctype.WithdrawMsg{
							Content:   content,
							MsgID:     request.Msg.WithdrawMsg.MsgID,
							OriginMsg: &msgModel.Msg,
						},
					},
				})

				// 把原消息置空
			}

			// 先入库
			msgID := InsertMsgByChat(svcCtx.DB, request.RevUserID, req.UserID, request.Msg)
			// 判断目标用户在不在线 给发送双方都要发消息
			SendMsgByUser(svcCtx, request.RevUserID, req.UserID, request.Msg, msgID)
		}

	}
}

type ChatRequest struct {
	RevUserID uint      `json:"revUserID"` // 给谁发
	Msg       ctype.Msg `json:"msg"`
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
		SendTipErrMsg(sendUser.Conn, "消息保存失败")
	}
	return chatModel.ID
}

// SendMsgByUser 发消息 给谁发 谁发的
func SendMsgByUser(svcCtx *svc.ServiceContext, revUserID uint, sendUserID uint, msg ctype.Msg, msgID uint) {

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
		revUser.Conn.WriteMessage(websocket.TextMessage, byteData)
		return
	}

	// 在线的情况下, 我是可以拿到对方的在线信息的
	// 对方不在线的情况下, 我只能通过调用用户服务rpc方法来获取

	// 无论如何都要给发送者回传消息

	if !ok1 {
		UserBaseInfo, err := svcCtx.UserRpc.UserBaseInfo(context.Background(), &user_rpc.UserBaseInfoRequest{
			UserId: uint32(revUserID),
		})

		if err != nil {
			logx.Error(err)
			return
		}

		resp.RevUser = ctype.UserInfo{
			ID:       revUserID,
			NickName: UserBaseInfo.NickName,
			Avatar:   UserBaseInfo.Avatar,
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
	sendUser.Conn.WriteMessage(websocket.TextMessage, byteData)

	//if ok1 && ok2 && sendUserID == revUserID{ // 自己加的判断是否为自己
	//	resp.IsMe = true
	//	byteData, _ = json.Marshal(resp)
	//	revUser.Conn.WriteMessage(websocket.TextMessage, byteData)
	//} else
	if ok1 {
		// 接受者在线
		resp.IsMe = false
		byteData, _ = json.Marshal(resp)
		revUser.Conn.WriteMessage(websocket.TextMessage, byteData)

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
