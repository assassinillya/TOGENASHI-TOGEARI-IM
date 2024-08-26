package logic

import (
	"context"
	"im_server/common/list_query"
	"im_server/common/models"
	"im_server/im_group/group_models"
	"im_server/im_user/user_rpc/types/user_rpc"

	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupValidListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupValidListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupValidListLogic {
	return &GroupValidListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupValidListLogic) GroupValidList(req *types.GroupValidListRequest) (resp *types.GroupValidListResponse, err error) {
	// 群验证列表  自己得是群管理员或者是群主
	var groupIDList []uint // 我管理的群
	l.svcCtx.DB.Model(group_models.GroupMemberModel{}).
		Where("user_id = ? and (role = 1 or role = 2)", req.UserID).
		Select("group_id").Scan(&groupIDList)

	// 先去查自己管理了哪些群，然后去找这些群的验证表

	groups, count, err := list_query.ListQuery(l.svcCtx.DB, group_models.GroupVerifyModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Preloads: []string{"GroupModel"},
		Where:    l.svcCtx.DB.Where("group_id in ? or user_id = ?", groupIDList, req.UserID),
	})

	var userIDList []uint32
	for _, group := range groups {
		userIDList = append(userIDList, uint32(group.UserID))
	}

	// 是谁想加群?
	userList, err1 := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})

	resp = new(types.GroupValidListResponse)
	resp.Count = int(count)
	for _, group := range groups {
		info := types.GroupValidInfoResponse{
			ID:                 group.ID,
			GroupID:            group.GroupID,
			UserID:             group.UserID,
			Status:             group.Status,
			AdditionalMessages: group.AdditionalMessages,
			Title:              group.GroupModel.Title,
			CreatedAt:          group.CreatedAt.String(),
			Type:               group.Type,
		}
		if group.VerificationQuestion != nil {
			info.VerificationQuestion = &types.VerificationQuestion{
				Problem1: group.VerificationQuestion.Problem1,
				Problem2: group.VerificationQuestion.Problem2,
				Problem3: group.VerificationQuestion.Problem3,
				Answer1:  group.VerificationQuestion.Answer1,
				Answer2:  group.VerificationQuestion.Answer2,
				Answer3:  group.VerificationQuestion.Answer3,
			}
		}
		if err1 == nil {
			info.UserNickname = userList.UserInfo[uint32(info.UserID)].NickName
			info.UserAvatar = userList.UserInfo[uint32(info.UserID)].Avatar
		}

		resp.List = append(resp.List, info)
	}

	return
}
