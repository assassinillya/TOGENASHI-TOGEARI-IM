package logic

import (
	"context"
	"errors"
	"fmt"
	"im_server/common/list_query"
	"im_server/common/models"
	"im_server/im_chat/chat_models"
	"im_server/im_user/user_rpc/types/user_rpc"

	"im_server/im_chat/chat_api/internal/svc"
	"im_server/im_chat/chat_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatSessionLogic {
	return &ChatSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Data struct {
	SU         uint   `gorm:"column:sU"`
	RU         uint   `gorm:"column:rU"`
	MaxDate    string `gorm:"column:maxDate"`
	MaxPreview string `gorm:"column:maxPreview"`
	IsTop      bool   `gorm:"column:isTop"`
}

func (l *ChatSessionLogic) ChatSession(req *types.ChatSessionRequest) (resp *types.ChatSessionResponse, err error) {
	column := fmt.Sprintf("if ((select 1 from top_user_models where user_id = %d and (top_user_id = sU or top_user_id = rU)limit 1), 1, 0) as isTop", req.UserID)
	chatList, count, _ := list_query.ListQuery(l.svcCtx.DB, Data{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "isTop desc, maxDate desc",
		},
		Table: func() (string, any) {
			return "(?) as u", l.svcCtx.DB.Model(&chat_models.ChatModel{}).
				Select("least(send_user_id, rev_user_id) as sU",
					" greatest(send_user_id, rev_user_id) as rU",
					"max(created_at) as maxDate",
					fmt.Sprintf("(select msg_preview from chat_models where ((send_user_id = sU and rev_user_id = rU) or (send_user_id = rU and rev_user_id = sU)) and  id not in (select chat_id from user_chat_delete_models where user_id = %d)  order by created_at desc limit 1 ) as maxPreview", req.UserID),
					column).
				Where("(send_user_id = ? or rev_user_id = ?) and id not in (select chat_id from user_chat_delete_models where user_id = ?)",
					req.UserID, req.UserID, req.UserID).
				Group("least(send_user_id, rev_user_id)").
				Group("greatest(send_user_id, rev_user_id)")
		},
	})

	var userIDList []uint32
	for _, data := range chatList {
		if data.RU != req.UserID {
			userIDList = append(userIDList, uint32(data.RU))
		}
		if data.SU != req.UserID {
			userIDList = append(userIDList, uint32(data.SU))
		}
		if data.SU == req.UserID && req.UserID == data.RU {
			// 自己和自己聊天
			userIDList = append(userIDList, uint32(req.UserID))
		}
	}
	// 使用客户端拦截器实现
	//fmt.Println(l.ctx.Value("clientIP").(string),l.ctx.Value("userID").(string))
	//clientIPStr:= l.ctx.Value("clientIP").(string)
	//userIDStr:= l.ctx.Value("userID").(string)
	//fmt.Println(clientIPStr,userIDStr)
	//md:= metadata.New(map[string]string{"clientIP":clientIPStr, "userID": userIDStr})
	//ctx := metadata.NewOutgoingContext(context.Background(), md)
	response, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})

	if err != nil {
		logx.Error(err)
		return nil, errors.New("用户服务错误")
	}

	var list = make([]types.ChatSession, 0)
	for _, data := range chatList {
		s := types.ChatSession{
			CreatedAt:  data.MaxDate,
			MsgPreview: data.MaxPreview,
			IsTop:      data.IsTop,
		}
		if data.RU != req.UserID {
			s.UserID = data.RU
			fmt.Println(uint32(s.UserID))
			s.Avatar = response.UserInfo[uint32(s.UserID)].Avatar
			s.Nickname = response.UserInfo[uint32(s.UserID)].NickName

		}
		if data.SU != req.UserID {
			s.UserID = data.SU
			s.Avatar = response.UserInfo[uint32(s.UserID)].Avatar
			s.Nickname = response.UserInfo[uint32(s.UserID)].NickName
		}
		if data.SU == req.UserID && req.UserID == data.RU {
			s.UserID = data.SU
			s.Avatar = response.UserInfo[uint32(s.UserID)].Avatar
			s.Nickname = response.UserInfo[uint32(s.UserID)].NickName
		}
		list = append(list, s)
	}

	return &types.ChatSessionResponse{
		List:  list,
		Count: count,
	}, nil
}
