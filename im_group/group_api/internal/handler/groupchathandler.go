package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"im_server/common/models/ctype"
	"im_server/common/response"
	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"
	"im_server/im_group/group_models"
	"im_server/im_user/user_rpc/types/user_rpc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type UserWsInfo struct {
	UserInfo    ctype.UserInfo             // 用户信息
	WsClientMap map[string]*websocket.Conn // 这个用户管理的所有ws客户端
}

var UserOnlineWsMap = map[uint]*UserWsInfo{}

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
			// 查在线的用户列表
			userOnlineIDList := getOnlineUserIDList()
			// 查这个群的成员 并且在线
			var groupMemberOnlineIDList []uint
			svcCtx.DB.Model(group_models.GroupMemberModel{}).
				Where("group_id = ? and user_id in ?", "", userOnlineIDList).
				Select("user_id").Scan(&groupMemberOnlineIDList)
			// 遍历这个用户列表，去找ws的客户端
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

func getOnlineUserIDList() (userOnlineIDList []uint) {
	for u, _ := range UserOnlineWsMap {
		userOnlineIDList = append(userOnlineIDList, u)
	}
	return
}
