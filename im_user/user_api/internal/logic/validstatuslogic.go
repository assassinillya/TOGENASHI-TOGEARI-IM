package logic

import (
	"context"
	"encoding/json"
	"errors"
	"im_server/common/models/ctype"
	"im_server/im_chat/chat_rpc/chat"
	"im_server/im_user/user_models"

	"im_server/im_user/user_api/internal/svc"
	"im_server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidStatusLogic {
	return &ValidStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidStatusLogic) ValidStatus(req *types.FriendValidStatusRequest) (resp *types.FriendValidResponse, err error) {

	var friendVerify user_models.FriendVerifyModel
	// 要操作状态, 前提我自己是接收方
	err = l.svcCtx.DB.Take(&friendVerify, "id = ? and rev_user_id = ?", req.VerifyID, req.UserID).Error
	if err != nil {
		return nil, errors.New("验证记录不存在")
	}

	if friendVerify.RevStatus != 0 {
		return nil, errors.New("验证记录已处理")
	}

	switch req.Status {
	case 1: // 同意
		friendVerify.RevStatus = 1
		// 往好友表里面加
		l.svcCtx.DB.Create(&user_models.FriendModel{
			RevUserID:  friendVerify.RevUserID,
			SendUserID: friendVerify.SendUserID,
		})

		msg := ctype.Msg{
			Type: ctype.TextMsgType,
			TextMsg: &ctype.TextMsg{
				Content: "我们已经是好友了，开始聊天吧！",
			},
		}
		byteData, _ := json.Marshal(msg)

		// 给对方发个消息
		_, err = l.svcCtx.ChatRpc.UserChat(l.ctx, &chat.UserChatRequest{
			SendUserId: uint32(friendVerify.SendUserID),
			RevUserId:  uint32(friendVerify.RevUserID),
			Msg:        byteData,
			SystemMsg:  nil,
		})

		if err != nil {
			logx.Error(err)
		}
	case 2: // 拒绝
		friendVerify.RevStatus = 2
	case 3: // 忽略
		friendVerify.RevStatus = 3
	case 4: // 删除
		// 一条验证记录, 是给两个人看的
		l.svcCtx.DB.Delete(&friendVerify)
		return nil, nil
	}
	l.svcCtx.DB.Save(&friendVerify)

	return
}
