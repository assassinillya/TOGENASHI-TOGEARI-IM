package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"im_server/im_user/user_api/internal/svc"
	"im_server/im_user/user_api/internal/types"
	"im_server/im_user/user_models"
	"im_server/im_user/user_rpc/types/user_rpc"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {

	res, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoRequest{
		UserId: uint32(req.UserID),
	})

	if err != nil {
		return nil, err
	}

	var user user_models.UserModel
	err = json.Unmarshal(res.Data, &user)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("数据错误")
	}

	resp = &types.UserInfoResponse{
		UserID:        user.ID,
		NickName:      user.Nickname,
		Abstract:      user.Abstract,
		Avatar:        user.Avatar,
		RecallMessage: user.UserConfModel.RecallMessage,
		FriendOnline:  user.UserConfModel.FriendOnline,
		Sound:         user.UserConfModel.Sound,
		SecureLink:    user.UserConfModel.SecureLink,
		SavePwd:       user.UserConfModel.SavePwd,
		SearchUser:    user.UserConfModel.SearchUser,
		Verification:  user.UserConfModel.Verification,
	}

	if user.UserConfModel.VerificationQuestion != nil {
		resp.VerificationQuestion = &types.VerificationQuestion{
			Problem1: user.UserConfModel.VerificationQuestion.Problem1,
			Problem2: user.UserConfModel.VerificationQuestion.Problem2,
			Problem3: user.UserConfModel.VerificationQuestion.Problem3,
			Answer1:  user.UserConfModel.VerificationQuestion.Answer1,
			Answer2:  user.UserConfModel.VerificationQuestion.Answer2,
			Answer3:  user.UserConfModel.VerificationQuestion.Answer3,
		}
	}

	return resp, nil
}
