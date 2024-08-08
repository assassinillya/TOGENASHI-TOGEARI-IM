package logic

import (
	"context"
	"im_server/im_chat/chat_models"

	"im_server/im_chat/chat_api/internal/svc"
	"im_server/im_chat/chat_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatDeleteLogic {
	return &ChatDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatDeleteLogic) ChatDelete(req *types.ChatDeleteRequest) (resp *types.ChatDeleteResponse, err error) {

	var chatList []chat_models.ChatModel
	l.svcCtx.DB.Find(&chatList, req.IdList)

	var deleteChatList []chat_models.UserChatDeleteModel

	var userDeleteChatList []chat_models.UserChatDeleteModel
	l.svcCtx.DB.Find(&userDeleteChatList, req.IdList)
	chatDeleteMap := map[uint]struct{}{}
	for _, model := range userDeleteChatList {
		chatDeleteMap[model.ChatID] = struct{}{}
	}

	if len(chatList) > 0 {
		for _, model := range chatList {
			// 不是自己的聊天记录
			if !(model.SendUserID == req.UserID || model.RevUserID == req.UserID) {
				logx.Info(model.ID, " 不是自己的聊天记录")
				continue
			}
			// 已经删过的聊天记录

			_, ok := chatDeleteMap[model.ID]
			if ok {
				logx.Info(model.ID, " 已经被删除过了")
				continue
			}

			deleteChatList = append(deleteChatList, chat_models.UserChatDeleteModel{
				UserID: req.UserID,
				ChatID: model.ID,
			})
		}
	}

	if len(deleteChatList) > 0 {
		l.svcCtx.DB.Create(&deleteChatList)
	}

	logx.Infof("已经删除了 %d 条记录", len(deleteChatList))

	return
}
