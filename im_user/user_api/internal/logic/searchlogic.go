package logic

import (
	"context"
	"im_server/common/list_quary"
	"im_server/common/models"
	"im_server/im_user/user_models"

	"im_server/im_user/user_api/internal/svc"
	"im_server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	// 先找所有的用户
	friends, count, _ := list_quary.ListQuery(l.svcCtx.DB, user_models.UserConfModel{}, list_quary.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Preload: []string{"UserModel"},
		Where:   l.svcCtx.DB.Where("search_user <> 0 or (search_user = 1 and user_id = ?)", req.Key),
	})
	list := make([]types.SearchInfo, 0)
	for _, friend := range friends {
		list = append(list, types.SearchInfo{
			UserID:   friend.UserID,
			NickName: friend.UserModel.Nickname,
			Abstract: friend.UserModel.Abstract,
			Avatar:   friend.UserModel.Avatar,
		})

	}
	return &types.SearchResponse{
		Count: count,
		List:  list,
	}, nil
}
