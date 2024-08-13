package handler

import (
	"context"
	"errors"
	"fmt"
	"im_server/common/response"
	"im_server/im_file/file_api/internal/logic"
	"im_server/im_file/file_api/internal/svc"
	"im_server/im_file/file_api/internal/types"
	"im_server/im_user/user_rpc/types/user_rpc"
	"im_server/utils"
	"im_server/utils/random"
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
		dir, err := os.ReadDir(dirPath)
		if err != nil {
			os.MkdirAll(dirPath, 0666)
		}

		fileName := fileHead.Filename
		filePath := path.Join(dirPath, fileHead.Filename)

		imageData, _ := io.ReadAll(file)
		l := logic.NewFileLogic(r.Context(), svcCtx)
		resp, err := l.File(&req)
		resp.Src = "/" + filePath
		if utils.InDir(dir, fileName) {
			// 文件重名
			prefix := utils.GetFilePrefix(fileName)
			newPath := fmt.Sprintf("%s_%s.%s", prefix, random.RandStr(4), suffix)
			filePath = path.Join(dirPath, newPath)
		}

		err = os.WriteFile(filePath, imageData, 0666)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp.Src = "/" + filePath
		response.Response(r, w, resp, err)
	}
}
