package logic

import (
	"context"
	"encoding/json"
	"im_server/common/models/ctype"
	"im_server/im_chat/chat_models"

	"im_server/im_chat/chat_rpc/internal/svc"
	"im_server/im_chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCreateLogic) UserCreate(in *chat_rpc.UserChatRequest) (*chat_rpc.UserChatResponse, error) {

	var msg ctype.Msg
	err := json.Unmarshal(in.Msg, &msg)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	var systemMsg *ctype.SystemMsg
	err = json.Unmarshal(in.SystemMsg, &systemMsg)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	err = l.svcCtx.DB.Create(&chat_models.ChatModel{
		SendUserID: uint(in.SendUserId),
		RevUserID:  uint(in.RevUserId),
		MsgType:    msg.Type,
		MsgPreView: "", //todo 写个方法获取
		Msg:        msg,
		SystemMsg:  systemMsg,
	}).Error

	if err != nil {
		logx.Error(err)
		return nil, err
	}

	return &chat_rpc.UserChatResponse{}, nil
}
