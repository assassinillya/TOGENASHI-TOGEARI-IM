package handler

import (
	"im_server/common/response"
	"im_server/im_group/group_api/internal/logic"
	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func groupCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupCreateRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGroupCreateLogic(r.Context(), svcCtx)
		resp, err := l.GroupCreate(&req)
		response.Response(r, w, resp, err)

	}
}
