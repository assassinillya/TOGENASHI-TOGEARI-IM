package logic

import (
	"context"
	"errors"
	"fmt"
	"im_server/im_group/group_models"
	"time"

	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupProhibitionUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupProhibitionUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupProhibitionUpdateLogic {
	return &GroupProhibitionUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupProhibitionUpdateLogic) GroupProhibitionUpdate(req *types.GroupProhibitionRequest) (resp *types.GroupProhibitionResponse, err error) {
	// 调用接口用户
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("用户服务错误")
	}
	if member.Role == 3 {
		return nil, errors.New("当前用户并非管理员")
	}

	// 被禁言用户
	var member1 group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member1, "group_id = ? and user_id = ?", req.GroupID, req.MemberID).Error
	if err != nil {
		return nil, errors.New("目标用户不是群成员")
	}

	if member.UserID == member1.UserID {
		return nil, errors.New("无法禁言自己")
	}

	if member1.Role == 1 {
		return nil, errors.New("无法禁言群主")
	}

	l.svcCtx.DB.Model(&member1).Update("prohibition_time", req.ProhibitionTime)

	// 利用redis的过期时间去做这个禁言时间
	key := fmt.Sprintf("prohibition__%d", member1.ID)
	if req.ProhibitionTime != nil {
		// 给redis设置一个key，过期时间是
		l.svcCtx.Redis.Set(key, "1", time.Duration(*req.ProhibitionTime)*time.Minute)
	} else {
		l.svcCtx.Redis.Del(key)
	}

	return
}
