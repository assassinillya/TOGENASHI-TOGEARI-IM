package handler

import (
	"im_server/common/response"
	"im_server/im_file/file_api/internal/svc"
	"im_server/im_file/file_api/internal/types"
	"net/http"
	"os"
	"path"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ImageShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImageShowRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		filePath := path.Join("uploads", req.ImageType, req.ImageName)

		byteData, err := os.ReadFile(filePath)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		w.Write(byteData)

	}
}
