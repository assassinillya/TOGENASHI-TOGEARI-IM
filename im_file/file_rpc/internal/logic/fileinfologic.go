package logic

import (
	"context"
	"errors"
	"im_server/im_file/file_model"
	"strings"

	"im_server/im_file/file_rpc/internal/svc"
	"im_server/im_file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileInfoLogic {
	return &FileInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FileInfoLogic) FileInfo(in *file_rpc.FileInfoRequest) (*file_rpc.FileInfoResponse, error) {
	var fileModel file_model.FileModel
	err := l.svcCtx.DB.Take(&fileModel, "uid = ?", in.FileId).Error
	if err != nil {
		return nil, errors.New("文件不存在")
	}

	var tp string
	nameList := strings.Split(fileModel.FileName, ".")
	if len(nameList) > 1 {
		tp = nameList[len(nameList)-1]
	}

	return &file_rpc.FileInfoResponse{
		FileName: fileModel.FileName,
		FileHash: fileModel.Hash,
		FileSize: fileModel.Size,
		FileType: tp,
	}, nil
}
