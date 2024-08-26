package logic

import (
	"context"
	"errors"
	"fmt"
	"im_server/common/list_query"
	"im_server/common/models"
	"im_server/common/models/ctype"
	"im_server/im_group/group_models"
	"im_server/im_user/user_rpc/types/user_rpc"

	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberLogic {
	return &GroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Data struct {
	GroupID        uint   `gorm:"column:group_id"`
	UserID         uint   `gorm:"column:user_id"`
	Role           int8   `gorm:"column:role"`
	CreatedAt      string `gorm:"column:created_at"`
	MemberNickname string `gorm:"column:member_nickname"`
	NewMsgDate     string `gorm:"column:new_msg_date"`
}

func (l *GroupMemberLogic) GroupMember(req *types.GroupMemberRequest) (resp *types.GroupMemberResponse, err error) {

	switch req.Sort {
	case "new_msg_date desc", "new_msg_date asc": // 按照最新发言排序
	case "role asc": //  按照角色升序
	case "created_at desc", "created_at asc": // 按照进群时间排
	default:
		return nil, errors.New("不支持的排序模式")
	}

	column := fmt.Sprintf(fmt.Sprintf("(select group_msg_models.created_at from group_msg_models where group_member_models.group_id = %d  and group_msg_models.send_user_id = user_id) as new_msg_date", req.ID))

	memberList, count, _ := list_query.ListQuery(l.svcCtx.DB, Data{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  req.Sort,
		},
		Table: func() (string, any) {
			return "(?) as u", l.svcCtx.DB.Model(&group_models.GroupMemberModel{GroupID: req.ID}).
				Select("group_id",
					"user_id",
					"role",
					"created_at",
					"member_nickname",
					column)
		},
	})
	// 把用户id全部拿出来
	var userIDList []uint32
	for _, data := range memberList {
		userIDList = append(userIDList, uint32(data.UserID))
	}

	var userInfoMap = map[uint]ctype.UserInfo{}
	userListResponse, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})
	// 关于降级
	if err == nil {
		for u, info := range userListResponse.UserInfo {
			userInfoMap[uint(u)] = ctype.UserInfo{
				ID:       uint(u),
				NickName: info.NickName,
				Avatar:   info.Avatar,
			}
		}
	} else {
		logx.Error(err)
	}

	var userOnlineMap = map[uint]bool{}
	userOnlineResponse, err := l.svcCtx.UserRpc.UserOnlineList(l.ctx, &user_rpc.UserOnlineListRequest{})
	if err == nil {
		for _, u := range userOnlineResponse.UserIdList {
			userOnlineMap[uint(u)] = true
		}
	} else {
		logx.Error(err)
	}

	resp = new(types.GroupMemberResponse)

	for _, data := range memberList {
		resp.List = append(resp.List, types.GroupMemberInfo{
			UserID:         data.UserID,
			UserNickname:   userInfoMap[data.UserID].NickName,
			Avatar:         userInfoMap[data.UserID].Avatar,
			IsOnline:       userOnlineMap[data.UserID],
			Role:           data.Role,
			MemberNickname: data.MemberNickname,
			CreatedAt:      data.CreatedAt,
			NewMsgDate:     data.NewMsgDate,
		})
	}
	resp.Count = int(count)
	return
}
