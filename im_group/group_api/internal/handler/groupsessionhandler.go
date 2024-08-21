package handler

import (
	"im_server/common/response"
	"im_server/im_group/group_api/internal/logic"
	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func groupSessionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupSessionRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGroupSessionLogic(r.Context(), svcCtx)
		resp, err := l.GroupSession(&req)
		response.Response(r, w, resp, err)

	}
}
