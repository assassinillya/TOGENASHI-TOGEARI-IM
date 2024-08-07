package logic

import (
	"context"
	"im_server/common/list_quary"
	"im_server/common/models"
	"im_server/common/models/ctype"
	"im_server/im_chat/chat_api/internal/svc"
	"im_server/im_chat/chat_api/internal/types"
	"im_server/im_chat/chat_models"
	"im_server/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatHistoryLogic {
	return &ChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type ChatHistory struct {
	ID        uint             `json:"id"`
	SendUser  UserInfo         `json:"sendUser"`
	RevUser   UserInfo         `json:"revUser"`
	IsMe      bool             `json:"isMe"`      // 哪条消息是我发的
	CreatedAt string           `json:"createdAt"` // 消息时间
	Msg       ctype.Msg        `json:"msg"`       // 消息内容
	SystemMsg *ctype.SystemMsg `json:"systemMsg"` // 系统消息
}

type ChatHistoryResponse struct {
	List  []ChatHistory `json:"list"`
	Count int64         `json:"count"`
}

func (l *ChatHistoryLogic) ChatHistory(req *types.ChatHistoryRequest) (resp *ChatHistoryResponse, err error) {

	chatList, count, _ := list_quary.ListQuery(l.svcCtx.DB, chat_models.ChatModel{}, list_quary.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Where: l.svcCtx.DB.Where("send_user_id = ? or rev_user_id = ?", req.UserID, req.UserID),
	})

	var userIDList []uint
	for _, model := range chatList {
		userIDList = append(userIDList, model.SendUserID)
		userIDList = append(userIDList, model.RevUserID)
	}

	// 去重
	userIDList = utils.DeduplicationList(userIDList)
	// 去调用户服务的rpc接口, 获取用户信息 { 用户id: {用户信息} }

	var list = make([]ChatHistory, 0)
	for _, model := range chatList {
		list = append(list, ChatHistory{
			ID:        model.ID,
			CreatedAt: model.CreatedAt.String(),
			Msg:       model.Msg,
			SystemMsg: model.SystemMsg,
		})
	}

	resp = &ChatHistoryResponse{
		List:  list,
		Count: count,
	}

	return
}
