package logic

import (
	"context"
	"im_server/common/list_query"
	"im_server/common/models"
	"im_server/im_user/user_models"

	"im_server/im_user/user_api/internal/svc"
	"im_server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListRequest) (resp *types.FriendListResponse, err error) {

	//var count int64
	//l.svcCtx.DB.
	//	Model(&user_models.FriendModel{}).
	//	Where("send_user_id = ? or rev_user_id = ?", req.UserID, req.UserID).
	//	Count(&count)
	//var friends []user_models.FriendModel
	//if req.Limit <= 0 {
	//	req.Limit = 10
	//}
	//if req.Page <= 0 {
	//	req.Page = 1
	//}
	//offset := (req.Page - 1) * req.Limit
	//
	//l.svcCtx.DB.
	//	Preload("SendUserModel").
	//	Preload("RevUserModel").
	//	Limit(req.Limit).
	//	Offset(offset).
	//	Find(&friends, "send_user_id = ? or rev_user_id = ?", req.UserID, req.UserID)

	friends, count, _ := list_query.ListQuery(l.svcCtx.DB, user_models.FriendModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Preload: []string{"SendUserModel", "RevUserModel"},
	})

	var list []types.FriendInfoResponse
	for _, friend := range friends {
		info := types.FriendInfoResponse{}
		if friend.SendUserID == req.UserID {
			// 发送方
			info = types.FriendInfoResponse{
				UserID:   friend.RevUserID,
				NickName: friend.RevUserModel.Nickname,
				Abstract: friend.RevUserModel.Abstract,
				Avatar:   friend.RevUserModel.Avatar,
				Notice:   friend.SendUserNotice,
			}
		}

		if friend.RevUserID == req.UserID {
			// 接收方
			info = types.FriendInfoResponse{
				UserID:   friend.SendUserID,
				NickName: friend.SendUserModel.Nickname,
				Abstract: friend.SendUserModel.Abstract,
				Avatar:   friend.SendUserModel.Avatar,
				Notice:   friend.RevUserNotice,
			}
		}
		list = append(list, info)

	}

	return &types.FriendListResponse{
		List:  list,
		Count: int(count),
	}, nil

}
