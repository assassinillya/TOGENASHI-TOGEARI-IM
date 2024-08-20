package logic

import (
	"context"
	"errors"
	"im_server/im_group/group_models"

	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberRemoveLogic {
	return &GroupMemberRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberRemoveLogic) GroupMemberRemove(req *types.GroupMemberRemoveRequest) (resp *types.GroupMemberRemoveResponse, err error) {
	// 能调用这个接口的只能是这个群的成员
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "user_id = ? and group_id = ?", req.UserID, req.ID).Error

	if err != nil {
		logx.Error(err)
		return nil, errors.New("并非该群成员")
	}
	if member.Role == 3 {
		return nil, errors.New("并非管理员")
	}

	var member1 group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member1, "user_id = ? and group_id = ?", req.MemberID, req.ID).Error

	if err != nil {
		logx.Error(err)
		return nil, errors.New("删除的成员并非该群成员")
	}
	// 群主可以踢所有人
	// 管理员只能踢普通成员
	if member.Role == 2 && member1.Role != 3 {
		return nil, errors.New("并非群主, 拼尽全力无法踢出管理员")
	}

	if member.UserID == member1.UserID {
		return nil, errors.New("无法踢出自己")
	}

	err = l.svcCtx.DB.Delete(member1).Error
	if err != nil {
		return nil, errors.New("移出群成员失败")
	}

	return
}