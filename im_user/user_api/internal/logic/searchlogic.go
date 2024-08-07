package logic

import (
	"context"
	"fmt"
	"im_server/common/list_query"
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
	users, count, _ := list_query.ListQuery(l.svcCtx.DB, user_models.UserConfModel{
		Online: req.Online,
	}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Preload: []string{"UserModel"},
		Joins:   "left join user_models um on um.id = user_conf_models.user_id",
		Where: l.svcCtx.DB.Where(
			"(user_conf_models.search_user <> 0 or user_conf_models.search_user is not null) and (user_conf_models.search_user = 1 and um.id = ?) or (user_conf_models.search_user = 2 and (um.id = ? or um.nickname like ?))",
			req.Key, req.Key, fmt.Sprintf("%%%s%%", req.Key)),
	})

	// 查自己是否为这个用户的好友
	var friend user_models.FriendModel
	friends := friend.Friends(l.svcCtx.DB, req.UserID)
	userMap := map[uint]bool{}
	for _, model := range friends {
		if model.SendUserID == req.UserID {
			userMap[model.RevUserID] = true
		} else {
			userMap[model.SendUserID] = true
		}
	}

	list := make([]types.SearchInfo, 0)
	for _, uc := range users {
		list = append(list, types.SearchInfo{
			UserID:   uc.UserID,
			NickName: uc.UserModel.Nickname,
			Abstract: uc.UserModel.Abstract,
			Avatar:   uc.UserModel.Avatar,
			IsFriend: userMap[uc.UserID],
		})

	}
	return &types.SearchResponse{
		Count: count,
		List:  list,
	}, nil
}
