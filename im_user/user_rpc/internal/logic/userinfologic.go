package logic

import (
	"context"
	"encoding/json"
	"errors"
	"im_server/im_user/user_models"

	"im_server/im_user/user_rpc/internal/svc"
	"im_server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user_rpc.UserInfoRequest) (*user_rpc.UserInfoResponse, error) {

	var user user_models.UserModel
	err := l.svcCtx.DB.Preload("UserConfModel").Take(&user, in.UserId).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	byteData, _ := json.Marshal(user)
	return &user_rpc.UserInfoResponse{Data: byteData}, nil
}
