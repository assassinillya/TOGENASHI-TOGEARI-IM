package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"im_server/common/response"
	"im_server/im_file/file_api/internal/logic"
	"im_server/im_file/file_api/internal/svc"
	"im_server/im_file/file_api/internal/types"
	"im_server/im_file/file_model"
	"im_server/im_user/user_rpc/types/user_rpc"
	"im_server/utils"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileRequest
		if err := httpx.ParseHeaders(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		file, fileHead, err := r.FormFile("file")
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 文件上传使用黑名单 exe php 无后缀文件
		// 文件后缀白名单
		nameList := strings.Split(fileHead.Filename, ".") // name.png 1.asily.png
		var suffix string
		if len(nameList) > 1 {
			suffix = nameList[len(nameList)-1]
		}

		if utils.InList(svcCtx.Config.BlackList, suffix) {
			response.Response(r, w, nil, errors.New("图片格式不正确"))
			return
		}

		fileData, _ := io.ReadAll(file)
		fileHash := utils.MD5(fileData)

		l := logic.NewFileLogic(r.Context(), svcCtx)
		resp, err := l.File(&req)

		var fileModel file_model.FileModel
		err = svcCtx.DB.Take(&fileModel, "hash = ?", fileHash).Error
		if err == nil {
			resp.Src = fileModel.WebPath()
			logx.Info("文件hash值重复, 文件名为: ", fileHead.Filename)
			response.Response(r, w, resp, err)
			return
		}

		// 文件重名
		// 在保存文件之前, 去读文件列表, 如果有重名的, 算一下他们两个的hash值, 若一样就不用保存
		// 当他们hash值不一样, 就把最新的文件重命名后再保存 {old_name}_xxx.{suffix}

		// 先去拿用户信息
		userResponse, err := svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{
			UserIdList: []uint32{uint32(req.UserID)},
		})

		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		dirName := fmt.Sprintf("%d_%s", req.UserID, userResponse.UserInfo[uint32(req.UserID)].NickName)

		dirPath := path.Join(svcCtx.Config.UploadDir, "file", dirName)
		_, err = os.ReadDir(dirPath)
		if err != nil {
			os.MkdirAll(dirPath, 0666)
		}

		// 文件信息入库
		UID := uuid.New()
		newFileModel := &file_model.FileModel{
			UserID:   req.UserID,
			FileName: fileHead.Filename,
			Size:     fileHead.Size,
			Hash:     fileHash,
			Uid:      UID,
		}
		newFileModel.Path = path.Join(dirPath, fmt.Sprintf("%s.%s", UID, suffix))
		err = os.WriteFile(newFileModel.Path, fileData, 0666)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		err = svcCtx.DB.Create(newFileModel).Error
		if err != nil {
			logx.Error(err)
			response.Response(r, w, resp, err)
			return
		}

		resp.Src = newFileModel.WebPath()
		response.Response(r, w, resp, err)
	}
}
