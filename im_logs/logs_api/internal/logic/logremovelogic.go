package logic

import (
	"context"
	"im_server/im_logs/logs_model"

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

	var logList []logs_model.LogModel
	l.svcCtx.DB.Find(&logList, req.IdList)
	if len(logList) > 0 {
		l.svcCtx.DB.Delete(&logList)
		l.svcCtx.ActionLogs.SetItem("删除日志条数", len(logList))
		logx.Infof("删除日志条数 %d", len(logList))
	}
	l.svcCtx.ActionLogs.Info("删除日志操作")
	l.svcCtx.ActionLogs.Save(l.ctx)
	return
}
