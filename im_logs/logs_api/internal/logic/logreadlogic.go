package logic

import (
	"context"

	"im_server/im_logs/logs_api/internal/svc"
	"im_server/im_logs/logs_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogReadLogic {
	return &LogReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogReadLogic) LogRead(req *types.LogReadRequest) (resp *types.LogReadResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
