package handler

import (
	"im_server/common/response"
	"im_server/im_auth/auth_api/internal/logic"
	"im_server/im_auth/auth_api/internal/svc"
	"net/http"
)

func authenticationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAuthenticationLogic(r.Context(), svcCtx) // 创建新的业务逻辑实例
		token := r.Header.Get("token")                         // 从请求头中获取 token
		resp, err := l.Authentication(token)                   // 调用业务逻辑进行认证
		response.Response(r, w, resp, err)                     // 返回响应
	}
}
