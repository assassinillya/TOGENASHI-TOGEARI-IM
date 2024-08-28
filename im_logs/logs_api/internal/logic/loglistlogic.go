package logic

import (
	"context"
	"im_server/common/list_query"
	"im_server/common/models"
	"im_server/im_logs/logs_model"

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

	logList, count, _ := list_query.ListQuery(l.svcCtx.DB, logs_model.LogModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Likes: []string{"ip", "user_nickname", "title"},
	})

	resp = new(types.LogListResponse)
	for _, model := range logList {
		resp.List = append(resp.List, types.LogInfoResponse{
			ID:           model.ID,
			CreatedAt:    model.CreatedAt.String(),
			LogType:      model.LogType,
			IP:           model.IP,
			Addr:         model.Addr,
			UserID:       model.UserID,
			UserNickname: model.UserNickname,
			UserAvatar:   model.UserAvatar,
			Level:        model.Level,
			Title:        model.Title,
			Content:      model.Content,
			Service:      model.Service,
			IsRead:       model.IsRead,
		})
	}
	resp.Count = int(count)
	return
}
