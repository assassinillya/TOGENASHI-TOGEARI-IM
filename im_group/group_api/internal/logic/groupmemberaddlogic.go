package logic

import (
	"context"
	"errors"
	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"
	"im_server/im_group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberAddLogic {
	return &GroupMemberAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberAddLogic) GroupMemberAdd(req *types.GroupMemberAddRequest) (resp *types.GroupMemberAddResponse, err error) {
	// 群成员邀请好友，得IsInvite=true
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Preload("GroupModel").Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("并非该群成员")
	}
	if member.Role == 3 {
		// 我是普通用户
		if !member.GroupModel.IsInvite {
			return nil, errors.New("本群未开放好友邀请入群功能")
		}
	}
	// 查一下哪些用户已经进群了
	var memberList []group_models.GroupMemberModel
	l.svcCtx.DB.Find(&memberList, "group_id = ? and user_id in ?", req.ID, req.MemberIDList)
	// 在前端不会出现邀请已经在群里的人进入群聊, 因此不用处理
	if len(memberList) > 0 {
		return nil, errors.New("已经有用户在群里了")
	}
	for _, memberID := range req.MemberIDList {
		memberList = append(memberList, group_models.GroupMemberModel{
			GroupID: req.ID,
			UserID:  memberID,
			Role:    3,
		})
	}

	// todo 没有处理被邀请人和邀请人的关系, 甚至可以邀请虚空人进入群聊

	err = l.svcCtx.DB.Create(&memberList).Error
	if err != nil {
		logx.Error(err)
	}

	return
}
