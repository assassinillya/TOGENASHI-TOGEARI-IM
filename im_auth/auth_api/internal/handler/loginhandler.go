package handler

import (
	"im_server/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im_server/im_auth/auth_api/internal/logic"
	"im_server/im_auth/auth_api/internal/svc"
	"im_server/im_auth/auth_api/internal/types"
)

func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		//if err != nil {
		//	httpx.ErrorCtx(r.Context(), w, err)
		//} else {
		//	httpx.OkJsonCtx(r.Context(), w, resp)
		//}
		response.Response(r, w, resp, err)
	}
}
