package handler

import (
	"errors"
	"im_server/common/response"
	"im_server/im_file/file_api/internal/logic"
	"im_server/im_file/file_api/internal/svc"
	"im_server/im_file/file_api/internal/types"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImageRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		file, fileHead, err := r.FormFile("image")
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		imageType := r.FormValue("imageType")
		if imageType == "" {
			response.Response(r, w, nil, errors.New("imageType不能为空"))
			return
		}

		byteData, _ := io.ReadAll(file)
		fileName := fileHead.Filename
		filePath := path.Join("uploads", imageType, fileName)
		err = os.WriteFile(filePath, byteData, 0666)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewImageLogic(r.Context(), svcCtx)
		resp, err := l.Image(&req)
		resp.Url = "/" + filePath
		response.Response(r, w, resp, err)
	}
}
