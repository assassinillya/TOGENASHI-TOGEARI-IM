package logic

import (
	"context"
	"errors"
	"fmt"
	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"
	"im_server/im_group/group_models"
	"im_server/im_user/user_rpc/types/user_rpc"
	"im_server/utils/set"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoLogic {
	return &GroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupInfoLogic) GroupInfo(req *types.GroupInfoRequest) (resp *types.GroupInfoResponse, err error) {

	fmt.Println(l.ctx.Value("clientIP"), l.ctx.Value("userID"))

	var groupModel group_models.GroupModel
	err = l.svcCtx.DB.Preload("MemberList").Take(&groupModel, req.ID).Error
	if err != nil {
		return nil, errors.New("群不存在")
	}
	// 谁能调这个接口 必须得是这个群的成员
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("该用户不是群成员")
	}

	resp = &types.GroupInfoResponse{
		GroupID:     groupModel.ID,
		Title:       groupModel.Title,
		Abstract:    groupModel.Abstract,
		MemberCount: len(groupModel.MemberList),
		Avatar:      groupModel.Avatar,
		Role:        member.Role,
	}
	// 查用户列表信息
	var userIDList []uint32
	var userAllIDList []uint32
	for _, model := range groupModel.MemberList {
		if model.Role == 1 || model.Role == 2 {
			userIDList = append(userIDList, uint32(model.UserID))
		}
		userAllIDList = append(userAllIDList, uint32(model.UserID))
	}

	userListResponse, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})
	if err != nil {
		logx.Error(err)
		return nil, errors.New("用户服务异常")
	}
	if userListResponse == nil {
		logx.Error("用户服务异常")
		return nil, errors.New("用户服务异常")
	}
	var creator types.UserInfo
	var adminList = make([]types.UserInfo, 0)

	// 算在线用户总数
	// 用户服务需要去写一个在线的用户列表的方法
	userOnlineResponse, err := l.svcCtx.UserRpc.UserOnlineList(l.ctx, &user_rpc.UserOnlineListRequest{})
	if err == nil {
		// 算群成员和总的在线人数成员，取交集
		slice := set.Intersect(userOnlineResponse.UserIdList, userAllIDList)
		resp.MemberOnlineCount = len(slice)
	}

	for _, model := range groupModel.MemberList {
		if model.Role == 3 {
			continue
		}
		userInfo := types.UserInfo{
			UserID:   model.UserID,
			Avatar:   userListResponse.UserInfo[uint32(model.UserID)].Avatar,
			Nickname: userListResponse.UserInfo[uint32(model.UserID)].NickName,
		}
		if model.Role == 1 {
			creator = userInfo
			continue
		}
		if model.Role == 2 {
			adminList = append(adminList, userInfo)
		}
	}
	resp.Creator = creator
	resp.AdminList = adminList

	return
}
