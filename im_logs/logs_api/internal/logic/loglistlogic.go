package logic

import (
	"context"

	"im_server/im_logs/logs_api/internal/svc"
	"im_server/im_logs/logs_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogListLogic {
	return &LogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogListLogic) LogList(req *types.LogListRequest) (resp *types.LogListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
