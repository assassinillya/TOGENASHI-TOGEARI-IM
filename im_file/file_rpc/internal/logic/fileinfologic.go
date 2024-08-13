package logic

import (
	"context"

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

	return &file_rpc.FileInfoResponse{}, nil
}
