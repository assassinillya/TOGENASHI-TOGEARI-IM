package logic

import (
	"context"
	"errors"
	"im_server/common/list_query"
	"im_server/common/models"
	"im_server/common/models/ctype"
	"im_server/im_chat/chat_api/internal/svc"
	"im_server/im_chat/chat_api/internal/types"
	"im_server/im_chat/chat_models"
	"im_server/im_user/user_rpc/types/user_rpc"
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

type ChatHistory struct {
	ID        uint             `json:"id"`
	SendUser  ctype.UserInfo   `json:"sendUser"`
	RevUser   ctype.UserInfo   `json:"revUser"`
	IsMe      bool             `json:"isMe"`      // 哪条消息是我发的
	CreatedAt string           `json:"createdAt"` // 消息时间
	Msg       ctype.Msg        `json:"msg"`       // 消息内容
	SystemMsg *ctype.SystemMsg `json:"systemMsg"` // 系统消息
}

type ChatHistoryResponse struct {
	List  []ChatHistory `json:"list"`
	Count int64         `json:"count"`
}

// ChatHistory 用户对话聊天记录
func (l *ChatHistoryLogic) ChatHistory(req *types.ChatHistoryRequest) (resp *ChatHistoryResponse, err error) {

	if req.UserID != req.FriendID {
		// 判断是否为好友
		res, err := l.svcCtx.UserRpc.IsFriend(l.ctx, &user_rpc.IsFriendRequest{
			User1: uint32(req.UserID),
			User2: uint32(req.FriendID),
		})
		if err != nil {
			logx.Error(err)
			return nil, err
		}

		if !res.IsFriend {
			return nil, errors.New("你们还不是好友")
		}
	}

	chatList, count, _ := list_query.ListQuery(l.svcCtx.DB, chat_models.ChatModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Where: l.svcCtx.DB.Where("((send_user_id = ? and rev_user_id = ?) or (rev_user_id = ? and send_user_id = ?)) and (id not in(select chat_id from user_chat_delete_models where user_id = ?))",
			req.UserID, req.FriendID, req.UserID, req.FriendID, req.UserID),
	})

	var userIDList []uint32
	for _, model := range chatList {
		userIDList = append(userIDList, uint32(model.SendUserID))
		userIDList = append(userIDList, uint32(model.RevUserID))
	}

	// 去重
	userIDList = utils.DeduplicationList(userIDList)
	// 去调用户服务的rpc接口, 获取用户信息 { 用户id: {用户信息} }
	response, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})

	if err != nil {
		logx.Error(err)
		return nil, errors.New("用户服务错误")
	}

	var list = make([]ChatHistory, 0)
	for _, model := range chatList {

		sendUser := ctype.UserInfo{
			ID:       model.SendUserID,
			NickName: response.UserInfo[uint32(model.SendUserID)].NickName,
			Avatar:   response.UserInfo[uint32(model.SendUserID)].Avatar,
		}

		revUser := ctype.UserInfo{
			ID:       model.RevUserID,
			NickName: response.UserInfo[uint32(model.RevUserID)].NickName,
			Avatar:   response.UserInfo[uint32(model.RevUserID)].Avatar,
		}
		info := ChatHistory{
			ID:        model.ID,
			CreatedAt: model.CreatedAt.String(),
			Msg:       model.Msg,
			SystemMsg: model.SystemMsg,
			SendUser:  sendUser,
			RevUser:   revUser,
		}

		if info.SendUser.ID == req.UserID {
			info.IsMe = true
		}

		list = append(list, info)
	}

	resp = &ChatHistoryResponse{
		List:  list,
		Count: count,
	}

	return
}
