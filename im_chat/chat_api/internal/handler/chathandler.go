package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"im_server/common/response"
	"im_server/im_chat/chat_api/internal/svc"
	"im_server/im_chat/chat_api/internal/types"
	"im_server/im_user/user_models"
	"im_server/im_user/user_rpc/types/user_rpc"
	"net/http"
)

type UserInfo struct {
	Nickname string `json:"nickName"`
	Avatar   string `json:"avatar"`
	UserID   uint   `json:"userID"`
}

type UserWsInfo struct {
	UserInfo UserInfo        //用户信息
	Conn     *websocket.Conn // 用户的ws连接对象
}

var UserWsMap = map[uint]UserWsInfo{}

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
			delete(UserWsMap, req.UserID)
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
			UserInfo: UserInfo{
				UserID:   req.UserID,
				Avatar:   userInfo.Avatar,
				Nickname: userInfo.Nickname,
			},
			Conn: conn,
		}

		UserWsMap[req.UserID] = userWsInfo

		if userInfo.UserConfModel.FriendOnline {
			// 如果好友开启了好友上线提醒

			// 查一下自己的好友是不是上线了

			// 查一下自己的好友列表, 返回用户id列表, 看看UserWsMap中是否存在, 如果存在就给自己发一个好友上线的消息
		}

		logx.Info(UserWsMap)

		defer conn.Close()
		for {
			// 消息类型，消息，错误
			_, p, err := conn.ReadMessage()
			if err != nil {
				// 用户断开聊天
				fmt.Println(err)
				break
			}
			fmt.Println(string(p))
			// 发送消息
			conn.WriteMessage(websocket.TextMessage, []byte("xxx"))
		}

	}
}
