package logic

import (
	"context"
	"errors"
	"im_server/im_group/group_models"

	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupRemoveLogic {
	return &GroupRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupRemoveLogic) GroupRemove(req *types.GroupRemoveRequest) (resp *types.GroupRemoveResponse, err error) {
	// 只能是群主才能调用
	var groupMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Preload("GroupModel").Take(&groupMember, "user_id = ? and group_id = ?", req.UserID, req.ID).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("群不存在或者用户不是该群成员")
	}
	if groupMember.Role != 1 {
		return nil, errors.New("只能群主才能解散群")
	}

	//关联的这个群的信息也要删掉
	var msgList []group_models.GroupMsgModel
	l.svcCtx.DB.Find(&msgList, "group_id = ?", req.ID).Delete(&msgList)
	//群成员也要删掉
	var memberList []group_models.GroupMemberModel
	l.svcCtx.DB.Find(&memberList, "group_id = ?", req.ID).Delete(&memberList)
	//群验证消息
	var vList []group_models.GroupVerifyModel
	l.svcCtx.DB.Find(&vList, "group_id = ?", req.ID).Delete(&vList)
	//群解散
	var group group_models.GroupModel
	l.svcCtx.DB.Take(&group, req.ID).Delete(&group)

	return
}
