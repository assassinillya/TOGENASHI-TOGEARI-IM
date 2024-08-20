package logic

import (
	"context"
	"im_server/common/list_query"
	"im_server/common/models"
	"im_server/im_user/user_models"
	"strconv"

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

	friends, count, _ := list_query.ListQuery(l.svcCtx.DB, user_models.FriendModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Where:    l.svcCtx.DB.Where("send_user_id = ? or rev_user_id = ?", req.UserID, req.UserID),
		Preloads: []string{"SendUserModel", "RevUserModel"},
	})

	// 查哪些用户在线
	onlineMap := l.svcCtx.Redis.HGetAll("online").Val()
	var onlineUserMap = map[uint]bool{}
	for key, _ := range onlineMap {
		val, err1 := strconv.Atoi(key)
		if err1 != nil {
			logx.Error(err1)
			continue
		}
		onlineUserMap[uint(val)] = true
	}

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
				IsOnline: onlineUserMap[friend.RevUserID], // 查看对方是否在线
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
				IsOnline: onlineUserMap[friend.SendUserID],
			}
		}
		list = append(list, info)

	}

	return &types.FriendListResponse{
		List:  list,
		Count: int(count),
	}, nil

}
