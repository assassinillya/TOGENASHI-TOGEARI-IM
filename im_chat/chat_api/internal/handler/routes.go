// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"im_server/im_chat/chat_api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/chat/history",
				Handler: chatHistoryHandler(serverCtx),
			},
		},
	)
}
