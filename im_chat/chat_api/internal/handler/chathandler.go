package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"im_server/common/models/ctype"
	"im_server/common/response"
	"im_server/im_chat/chat_api/internal/svc"
	"im_server/im_chat/chat_api/internal/types"
	"im_server/im_user/user_models"
	"im_server/im_user/user_rpc/types/user_rpc"
	"net/http"
	"time"
)

type UserWsInfo struct {
	UserInfo user_models.UserModel //用户信息
	Conn     *websocket.Conn       // 用户的ws连接对象
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

		UserWsMap[req.UserID] = userWsInfo
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

		for _, info := range friendsRes.FriendList {
			friend, ok := UserWsMap[uint(info.UserId)]
			if ok {
				text := fmt.Sprintf("好友%s上线了", UserWsMap[req.UserID].UserInfo.Nickname) // todo 这里修改为备注好点 备注(nickName)
				logx.Info(text)
				// 判断用户是否开了好友上线提示功能
				if friend.UserInfo.UserConfModel.FriendOnline {
					// 好友上线了
					friend.Conn.WriteMessage(websocket.TextMessage, []byte(text))
				}

			}
		}
		// 查一下自己的好友列表, 返回用户id列表, 看看UserWsMap中是否存在, 如果存在就给自己发一个好友上线的消息

		logx.Info(UserWsMap)

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
					SendTipErrMsg(conn, "用户服务, 请重试")
					return
				}

				if !isFriendRes.IsFriend {
					SendTipErrMsg(conn, "你们还不是好友")
				}
			}

			// 先入库

			// 判断目标用户在不在线
			SendMsgByUser(request.RevUserID, req.UserID, request.Msg)
		}

	}
}

type ChatRequest struct {
	RevUserID uint      `json:"revUserID"` // 给谁发
	Msg       ctype.Msg `json:"msg"`
}

type ChatResponse struct {
	RevUser   ctype.UserInfo `json:"revUser"`
	SendUser  ctype.UserInfo `json:"sendUser"`
	Msg       ctype.Msg      `json:"msg"`
	CreatedAt time.Time      `json:"createdAt"`
}

// InsertMsgByChat 消息入库
func InsertMsgByChat(revUserID uint, sendUserID uint, msg ctype.Msg) {

}

// SendMsgByUser 发消息 给谁发 谁发的
func SendMsgByUser(revUserID uint, sendUserID uint, msg ctype.Msg) {

	revUser, ok := UserWsMap[revUserID]
	if !ok {
		return
	}

	sendUser, ok := UserWsMap[sendUserID]
	if !ok {
		return
	}

	resp := ChatResponse{
		RevUser: ctype.UserInfo{
			ID:       revUserID,
			NickName: revUser.UserInfo.Nickname,
			Avatar:   revUser.UserInfo.Avatar,
		},
		SendUser: ctype.UserInfo{
			ID:       sendUserID,
			NickName: sendUser.UserInfo.Nickname,
			Avatar:   sendUser.UserInfo.Avatar,
		},
		Msg:       msg,
		CreatedAt: time.Now(),
	}
	byteData, _ := json.Marshal(resp)
	revUser.Conn.WriteMessage(websocket.TextMessage, byteData)
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
