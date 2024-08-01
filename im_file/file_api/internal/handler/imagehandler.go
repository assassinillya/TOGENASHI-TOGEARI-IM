package handler

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"im_server/common/response"
	"im_server/im_file/file_api/internal/logic"
	"im_server/im_file/file_api/internal/svc"
	"im_server/im_file/file_api/internal/types"
	"im_server/utils"
	"im_server/utils/random"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

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
		switch imageType {
		case "avatar", "group_avatar", "chat":
		default:
			response.Response(r, w, nil, errors.New("imageType只能为 avatar,group_avatar,chat"))
			return
		}

		// 文件大小限制
		mSize := float64(fileHead.Size) / float64(1024) / float64(1024)

		if mSize > svcCtx.Config.FileSize {
			response.Response(r, w, nil, fmt.Errorf("图片大小超过限制，最大只能上传%.2fMB大小的图片",
				svcCtx.Config.FileSize))
			return
		}

		// 文件后缀白名单
		nameList := strings.Split(fileHead.Filename, ".") // name.png 1.asily.png
		var suffix string
		if len(nameList) > 1 {
			suffix = nameList[len(nameList)-1]
		}

		if !utils.InList(svcCtx.Config.WhiteList, suffix) {
			response.Response(r, w, nil, errors.New("图片格式不正确"))
			return
		}

		// TODO 这里的文件上传逻辑很简陋, 这个对比hash值只能对比最初的重名文件
		// 文件重名
		// 在保存文件之前, 去读文件列表, 如果有重名的, 算一下他们两个的hash值, 若一样就不用保存
		// 当他们hash值不一样, 就把最新的文件重命名后再保存 {old_name}_xxx.{suffix}

		dirPath := path.Join(svcCtx.Config.UploadDir, imageType)
		dir, err := os.ReadDir(dirPath)
		if err != nil {
			os.MkdirAll(dirPath, 0666)
		}

		imageData, _ := io.ReadAll(file)
		fileName := fileHead.Filename
		filePath := path.Join(svcCtx.Config.UploadDir, imageType, fileName)
		l := logic.NewImageLogic(r.Context(), svcCtx)
		resp, err := l.Image(&req)
		resp.Url = "/" + filePath
		if InDir(dir, fileHead.Filename) {
			// 文件重名
			// 先读取文件列表, 查一下他们两个的hash值
			byteData, _ := os.ReadFile(filePath)
			oldFileHash := utils.MD5(byteData)
			newFileHash := utils.MD5(imageData)
			if oldFileHash == newFileHash {
				// 文件hash值一样, 不用保存
				logx.Info("文件hash值一样, 不用保存")
				response.Response(r, w, resp, nil)
				return
			}

			// 文件hash值不一样, 重命名
			prefix := utils.GetFilePrefix(fileName)

			newPath := fmt.Sprintf("%s_%s.%s", prefix, random.RandStr(4), suffix)
			filePath = path.Join(svcCtx.Config.UploadDir, imageType, newPath)

			//改了名字后还是重名 这个地方就得递归判断了

		}

		err = os.WriteFile(filePath, imageData, 0666)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp.Url = "/" + filePath
		response.Response(r, w, resp, err)
	}
}

func InDir(dir []os.DirEntry, file string) bool {
	for _, entry := range dir {
		if entry.Name() == file {
			return true
		}
	}
	return false
}
