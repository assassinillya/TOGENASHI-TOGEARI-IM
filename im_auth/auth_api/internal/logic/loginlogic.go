package logic

import (
	"context"

	"im_server/im_auth/auth_api/internal/svc"
	"im_server/im_auth/auth_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line

	return &types.LoginResponse{
		Code: 0,
		Msg:  "登录成功",
		Data: types.LoginInfo{
			Token: "123456",
		},
	}, nil
}
