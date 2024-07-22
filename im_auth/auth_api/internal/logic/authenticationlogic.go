package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"im_server/im_auth/auth_api/internal/svc"
)

type AuthenticationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	return &AuthenticationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthenticationLogic) Authentication() (resp string, err error) {
	// todo: add your logic here and delete this line

	return
}
