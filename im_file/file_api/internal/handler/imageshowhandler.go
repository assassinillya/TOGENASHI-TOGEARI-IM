package handler

import (
	"errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"im_server/common/response"
	"im_server/im_file/file_api/internal/svc"
	"im_server/im_file/file_api/internal/types"
	"im_server/im_file/file_model"
	"net/http"
	"os"
)

func ImageShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImageShowRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		var fileModel file_model.FileModel
		err := svcCtx.DB.Take(&fileModel, "uid = ?", req.ImageName).Error
		if err != nil {
			response.Response(r, w, nil, errors.New("文件不存在"))
			return
		}

		byteData, err := os.ReadFile(fileModel.Path)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		w.Write(byteData)

	}
}
