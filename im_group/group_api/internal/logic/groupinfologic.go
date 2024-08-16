package logic

import (
	"context"
	"errors"
	"im_server/im_group/group_models"
	"im_server/im_user/user_rpc/users"

	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"

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
	// 能调用这个接口的只能是这个群的成员

	var groupModel group_models.GroupModel
	err = l.svcCtx.DB.Preload("MemberList").Take(&groupModel, req.ID).Error
	if err != nil {
		return nil, errors.New("群不存在")
	}
	// 算在线用户总数

	resp = &types.GroupInfoResponse{
		GroupID:     groupModel.ID,
		Title:       groupModel.Title,
		Abstract:    groupModel.Abstract,
		MemberCount: len(groupModel.MemberList),
		Avatar:      groupModel.Avatar,
	}

	// 查用户列表信息
	var userIDList []uint32
	for _, model := range groupModel.MemberList {
		if model.Role == 3 {
			continue
		}
		userIDList = append(userIDList, uint32(model.UserID))
	}

	userListResponse, err := l.svcCtx.UserRpc.UserListInfo(context.Background(), &users.UserListInfoRequest{
		UserIdList: userIDList,
	})
	if err != nil {
		return
	}

	var creator types.UserInfo
	var adminList = make([]types.UserInfo, 0)
	for _, model := range groupModel.MemberList {
		if model.Role == 1 {
			creator = types.UserInfo{
				UserID:   model.UserID,
				Avatar:   userListResponse.UserInfo[uint32(model.UserID)].Avatar,
				Nickname: userListResponse.UserInfo[uint32(model.UserID)].NickName,
			}
		}
		if model.Role == 2 {
			cnt := types.UserInfo{
				UserID:   model.UserID,
				Avatar:   userListResponse.UserInfo[uint32(model.UserID)].Avatar,
				Nickname: userListResponse.UserInfo[uint32(model.UserID)].NickName,
			}
			adminList = append(adminList, cnt)
		}
	}
	resp.Creator = creator
	resp.AdminList = adminList

	// 查在线用户数量
	// 用户服务需要去写一个在线的用户列表的接口的方法

	return
}
