package logic

import (
	"context"
	"errors"
	"im_server/im_user/user_models"

	"im_server/im_user/user_api/internal/svc"
	"im_server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendDeleteLogic {
	return &FriendDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendDeleteLogic) FriendDelete(req *types.FriendDeleteRequest) (resp *types.FriendDeleteResponse, err error) {
	// 判断是否已经是好友了
	var friend user_models.FriendModel
	if !friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID) {
		return nil, errors.New("你们还不是好友")
	}

	l.svcCtx.DB.Delete(&friend)

	return
}
