// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"im_server/im_group/group_api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/group/friends",
				Handler: groupFriendsListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/group/group",
				Handler: groupCreateHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/group/group/:id",
				Handler: groupInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/group/group/:id",
				Handler: groupRemoveHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/group/groupUpdate",
				Handler: groupUpdateHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/group/member",
				Handler: groupMemberHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/group/member",
				Handler: groupMemberRemoveHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/group/member",
				Handler: groupMemberAddHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/group/member/nickname",
				Handler: groupMemberNicknameUpdateHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/group/member/role",
				Handler: groupMemberRoleUpdateHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/group/search",
				Handler: groupSearchHandler(serverCtx),
			},
		},
	)
}
