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
	"im_server/utils"
	"time"

	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupHistoryLogic {
	return &GroupHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// HistoryResponse 单条信息
type HistoryResponse struct {
	UserID       uint      `json:"userID"`
	UserNickname string    `json:"userNickname"`
	UserAvatar   string    `json:"userAvatar"`
	Msg          ctype.Msg `json:"msg"`
	ID           uint      `json:"id"`
	MsgType      int8      `json:"msgType"`
	CreatedAt    time.Time `json:"createdAt"`
	IsMe         bool      `json:"isMe"`
}

type HistoryListResponse struct {
	List  []HistoryResponse `json:"list"`
	Count int               `json:"count"`
}

func (l *GroupHistoryLogic) GroupHistory(req *types.GroupHistoryRequest) (resp *HistoryListResponse, err error) {
	// 谁能调这个接口 必须得是这个群的成员
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("并非群成员")
	}
	// 去查我删除了哪些聊天记录
	var msgIDList []uint
	err5 := l.svcCtx.DB.Model(group_models.GroupUserMsgDeleteModel{}).
		Where("group_id = ? and user_id = ?", req.ID, req.UserID).
		Select("msg_id").Scan(&msgIDList).Error
	fmt.Println(err5)

	var query = l.svcCtx.DB.Where("")
	// 不加上这个会导致在msgIDList为空的情况下什么也搜不到
	if len(msgIDList) > 0 {
		query.Where("id not in ?", msgIDList)
	}

	groupMsgList, count, err := list_query.ListQuery(l.svcCtx.DB, group_models.GroupMsgModel{GroupID: req.ID}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Where: l.svcCtx.DB.Where("id not in ? and group_id = ?", msgIDList, req.ID),
	})

	var userIDList []uint32
	for _, model := range groupMsgList {
		userIDList = append(userIDList, uint32(model.SendUserID))
	}
	userIDList = utils.DeduplicationList(userIDList)
	userListResponse, err1 := l.svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})

	var list = make([]HistoryResponse, 0)
	for _, model := range groupMsgList {
		info := HistoryResponse{
			UserID:    model.SendUserID,
			Msg:       model.Msg,
			ID:        model.ID,
			MsgType:   model.MsgType,
			CreatedAt: model.CreatedAt,
		}
		// 拿不到名称和头像也无所谓
		if err1 == nil {
			info.UserNickname = userListResponse.UserInfo[uint32(info.UserID)].NickName
			info.UserAvatar = userListResponse.UserInfo[uint32(info.UserID)].Avatar
		}
		if req.UserID == info.UserID {
			info.IsMe = true
		}
		list = append(list, info)
	}

	resp = new(HistoryListResponse)
	resp.List = list
	resp.Count = int(count)
	return
}
