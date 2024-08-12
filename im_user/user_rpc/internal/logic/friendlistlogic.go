package logic

import (
	"context"
	"im_server/common/list_query"
	"im_server/common/models"
	"im_server/im_user/user_models"

	"im_server/im_user/user_rpc/internal/svc"
	"im_server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *user_rpc.FriendListRequest) (*user_rpc.FriendListResponse, error) {
	friends, _, _ := list_query.ListQuery(l.svcCtx.DB, user_models.FriendModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Limit: -1, // 查全部
		},
		Where:    l.svcCtx.DB.Where("send_user_id = ? or rev_user_id = ?", in.User, in.User),
		Preloads: []string{"SendUserModel", "RevUserModel"},
	})

	var list []*user_rpc.FriendInfo
	for _, friend := range friends {
		info := user_rpc.FriendInfo{}
		if friend.SendUserID == uint(in.User) {
			// 发送方
			info = user_rpc.FriendInfo{
				UserId:   uint32(friend.RevUserID),
				NickName: friend.RevUserModel.Nickname,
				Avatar:   friend.RevUserModel.Avatar,
			}
		}

		if friend.RevUserID == uint(in.User) {
			// 接收方
			info = user_rpc.FriendInfo{
				UserId:   uint32(friend.SendUserID),
				NickName: friend.RevUserModel.Nickname,
				Avatar:   friend.RevUserModel.Avatar,
			}
		}
		list = append(list, &info)

	}

	return &user_rpc.FriendListResponse{
		FriendList: list,
	}, nil
}
