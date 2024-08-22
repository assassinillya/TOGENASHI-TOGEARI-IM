package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"im_server/common/models/ctype"
	"im_server/common/response"
	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"
	"im_server/im_group/group_models"
	"im_server/im_user/user_rpc/types/user_rpc"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type UserWsInfo struct {
	UserInfo    ctype.UserInfo             // 用户信息
	WsClientMap map[string]*websocket.Conn // 这个用户管理的所有ws客户端
}

var UserOnlineWsMap = map[uint]*UserWsInfo{}

type ChatRequest struct {
	GroupID uint      `json:"groupID"` // 群id
	Msg     ctype.Msg `json:"msg"`     // 消息
}

type ChatResponse struct {
	UserID       uint          `json:"userID"`
	UserNickname string        `json:"userNickname"`
	UserAvatar   string        `json:"userAvatar"`
	Msg          ctype.Msg     `json:"msg"`
	ID           uint          `json:"id"`
	MsgType      ctype.MsgType `json:"msgType"`
	CreatedAt    time.Time     `json:"createdAt"`
	IsMe         bool          `json:"isMe"`
}

func groupChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupChatRequest
		if err := httpx.ParseHeaders(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 升级为ws
		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// 鉴权 true表示放行，false表示拦截
				return true
			},
		}

		conn, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		// 用户可能会开多个客户端
		addr := conn.RemoteAddr().String()
		logx.Infof("用户建立ws连接 %s", addr)
		defer func() {
			conn.Close()

			userWsInfo, ok := UserOnlineWsMap[req.UserID]
			if ok {
				// 删除的退出的那个ws信息
				delete(userWsInfo.WsClientMap, addr)
			}
			if userWsInfo != nil && len(userWsInfo.WsClientMap) == 0 {
				// 全退完了
				delete(UserOnlineWsMap, req.UserID)
			}
		}()

		// 获取用户基本信息
		baseInfoResponse, err := svcCtx.UserRpc.UserBaseInfo(context.Background(), &user_rpc.UserBaseInfoRequest{
			UserId: uint32(req.UserID),
		})
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		userInfo := ctype.UserInfo{
			ID:       req.UserID,
			NickName: baseInfoResponse.NickName,
			Avatar:   baseInfoResponse.Avatar,
		}

		userWsInfo, ok := UserOnlineWsMap[req.UserID]
		if !ok {
			userWsInfo = &UserWsInfo{
				UserInfo: userInfo,
				WsClientMap: map[string]*websocket.Conn{
					addr: conn,
				},
			}
			// 代表这个用户第一次来
			UserOnlineWsMap[req.UserID] = userWsInfo
		}
		_, ok1 := userWsInfo.WsClientMap[addr]
		if !ok1 {
			// 代表这个用户二开及以上
			UserOnlineWsMap[req.UserID].WsClientMap[addr] = conn
		}

		for {
			// 消息类型，消息，错误
			_, p, err1 := conn.ReadMessage()
			if err1 != nil {
				// 用户断开聊天
				fmt.Println(err1)
				break
			}

			var request ChatRequest
			err = json.Unmarshal(p, &request)
			if err != nil {
				SendTipErrMsg(conn, "参数解析失败")
				continue
			}

			// 判断自己是不是这个群的成员
			var member group_models.GroupMemberModel
			err = svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", request.GroupID, req.UserID).Error
			if err != nil {
				SendTipErrMsg(conn, "你还不是该群的成员")
				continue
			}

			msgID := InsertMsg(svcCtx.DB, conn, member, request.Msg)

			// 遍历这个用户列表，去找ws的客户端
			sendGroupOnlineUserMsg(
				svcCtx.DB,
				request.GroupID,
				req.UserID,
				request.Msg,
				msgID,
			)

			// 查在线的用户列表
			userOnlineIDList := getOnlineUserIDList()
			// 查这个群的成员 并且在线
			var groupMemberOnlineIDList []uint
			svcCtx.DB.Model(group_models.GroupMemberModel{}).
				Where("group_id = ? and user_id in ?", request.GroupID, userOnlineIDList).
				Select("user_id").Scan(&groupMemberOnlineIDList)

			for _, u := range groupMemberOnlineIDList {
				wsUserInfo, ok2 := UserOnlineWsMap[u]
				if !ok2 {
					continue
				}
				for _, w2 := range wsUserInfo.WsClientMap {
					w2.WriteMessage(websocket.TextMessage, []byte(""))
				}
			}

			fmt.Println(string(p))
		}
	}
}

func InsertMsg(DB *gorm.DB, conn *websocket.Conn, member group_models.GroupMemberModel, msg ctype.Msg) uint {
	switch msg.Type {
	case ctype.WithdrawMsgType:
		fmt.Println("撤回消息是不入库的")
		return 0
	}
	groupModel := group_models.GroupMsgModel{
		GroupID:    member.GroupID,
		SendUserID: member.UserID,
		MsgType:    msg.Type,
		Msg:        msg,
	}
	groupModel.MsgPreView = groupModel.MsgPreviewMethod()
	err := DB.Create(&groupModel).Error
	if err != nil {
		logx.Error(err)
		SendTipErrMsg(conn, "消息保存失败")
		return 0
	}
	return groupModel.ID
}

func getOnlineUserIDList() (userOnlineIDList []uint) {
	for u, _ := range UserOnlineWsMap {
		userOnlineIDList = append(userOnlineIDList, u)
	}
	return
}

// SendTipErrMsg 发送错误提示的消息
func SendTipErrMsg(Conn *websocket.Conn, msg string) {
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
	Conn.WriteMessage(websocket.TextMessage, byteData)
}

// 给这个群的用户发消息
func sendGroupOnlineUserMsg(db *gorm.DB, groupID uint, userID uint, msg ctype.Msg, msgID uint) {

	// 查在线的用户列表
	userOnlineIDList := getOnlineUserIDList()
	// 查这个群的成员 并且在线
	var groupMemberOnlineIDList []uint
	db.Model(group_models.GroupMemberModel{}).
		Where("group_id = ? and user_id in ?", groupID, userOnlineIDList).
		Select("user_id").Scan(&groupMemberOnlineIDList)

	// 构造响应
	var chatResponse = ChatResponse{
		UserID:    userID,
		Msg:       msg,
		ID:        msgID,
		MsgType:   msg.Type,
		CreatedAt: time.Now(),
	}
	wsInfo, ok := UserOnlineWsMap[userID]
	if ok {
		chatResponse.UserNickname = wsInfo.UserInfo.NickName
		chatResponse.UserAvatar = wsInfo.UserInfo.Avatar
	}

	for _, u := range groupMemberOnlineIDList {
		wsUserInfo, ok2 := UserOnlineWsMap[u]
		if !ok2 {
			continue
		}
		// 判断isMe
		if wsUserInfo.UserInfo.ID == userID {
			chatResponse.IsMe = true
		}

		byteData, _ := json.Marshal(chatResponse)

		for _, w2 := range wsUserInfo.WsClientMap {
			w2.WriteMessage(websocket.TextMessage, byteData)
		}
	}
}
