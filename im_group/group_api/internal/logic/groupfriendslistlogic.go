package logic

import (
	"context"
	"im_server/im_group/group_models"
	"im_server/im_user/user_rpc/types/user_rpc"

	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupFriendsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupFriendsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupFriendsListLogic {
	return &GroupFriendsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupFriendsListLogic) GroupFriendsList(req *types.GroupFriendsListRequest) (resp *types.GroupFriendsListResponse, err error) {
	// 我的好友哪些在这个群里面

	// 需要去查我的好友列表
	friendResponse, err := l.svcCtx.UserRpc.FriendList(context.Background(), &user_rpc.FriendListRequest{
		User: uint32(req.UserID),
	})
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	// 这个群的群成员列表 组成一个map
	var memberList []group_models.GroupMemberModel
	l.svcCtx.DB.Find(&memberList, "group_id = ?", req.ID)
	var memberMap = map[uint]bool{}
	for _, model := range memberList {
		memberMap[model.UserID] = true
	}
	resp = new(types.GroupFriendsListResponse)

	for _, info := range friendResponse.FriendList {
		resp.List = append(resp.List, types.GroupFriendsResponse{
			UserID:    uint(info.UserId),
			Avatar:    info.Avatar,
			Nickname:  info.NickName,
			IsInGroup: memberMap[uint(info.UserId)],
		})
	}
	return
}
