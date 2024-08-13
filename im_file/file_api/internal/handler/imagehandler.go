package handler

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"im_server/common/response"
	"im_server/im_file/file_api/internal/logic"
	"im_server/im_file/file_api/internal/svc"
	"im_server/im_file/file_api/internal/types"
	"im_server/im_file/file_model"
	"im_server/utils"
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

		imageData, _ := io.ReadAll(file)
		imageHash := utils.MD5(imageData)

		var fileModel file_model.FileModel
		err = svcCtx.DB.Take(&fileModel, "hash = ?", imageHash).Error
		l := logic.NewImageLogic(r.Context(), svcCtx)
		resp, err := l.Image(&req)

		fmt.Println(fileModel, "and", err)
		if err == nil {
			// 找到了有一模一样的图片, 返回之前已有文件的WebPath
			resp.Url = fileModel.WebPath()
			logx.Info("文件hash值重复, 文件名为: ", fileHead.Filename)
			response.Response(r, w, resp, nil)
			return

		}

		// 拼接路径 /uploads/imageType/{uid}.{后缀}

		dirPath := path.Join(svcCtx.Config.UploadDir, imageType)
		_, err = os.ReadDir(dirPath)
		if err != nil {
			os.MkdirAll(dirPath, 0666)
		}

		fileName := fileHead.Filename
		UID := uuid.New()
		// 文件信息入库
		newFileModel := &file_model.FileModel{
			UserID:   req.UserID,
			FileName: fileName,
			Size:     fileHead.Size,
			Hash:     utils.MD5(imageData),
			Uid:      UID,
		}
		newFileModel.Path = path.Join(dirPath, fmt.Sprintf("%s.%s", UID, suffix))
		fmt.Println(newFileModel.Path)

		err = os.WriteFile(newFileModel.Path, imageData, 0666)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		err = svcCtx.DB.Create(newFileModel).Error
		if err != nil {
			logx.Error(err)
			response.Response(r, w, resp, err)
			return
		}

		resp.Url = newFileModel.WebPath()

		response.Response(r, w, resp, err)
	}
}
