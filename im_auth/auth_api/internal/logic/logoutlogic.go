package logic

import (
	"context"
	"errors"
	"fmt"
	"im_server/utils/jwts"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"im_server/im_auth/auth_api/internal/svc"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(token string) (resp string, err error) {
	if token == "" {
		err = errors.New("请传入token")
		return
	}

	payload, err := jwts.ParseToken(token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		err = errors.New("token错误")
		return
	}
	//设置过期时间为jwt的失效时间
	now := time.Now()
	expiration := payload.ExpiresAt.Time.Sub(now)

	key := fmt.Sprintf("logout_%s", token)
	l.svcCtx.Redis.SetNX(key, "", expiration)
	resp = "注销成功"

	return
}
