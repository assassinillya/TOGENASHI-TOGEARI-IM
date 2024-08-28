package logic

import (
	"context"

	"im_server/im_logs/logs_api/internal/svc"
	"im_server/im_logs/logs_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogRemoveLogic {
	return &LogRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogRemoveLogic) LogRemove(req *types.LogRemoveRequest) (resp *types.LogRemoveResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
