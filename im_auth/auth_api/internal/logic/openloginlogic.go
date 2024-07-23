package logic

import (
	"context"
	"errors"
	"im_server/im_auth/auth_models"
	"im_server/utils/open_login"
	"log"

	"im_server/im_auth/auth_api/internal/svc"
	"im_server/im_auth/auth_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Open_loginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpen_loginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Open_loginLogic {
	return &Open_loginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Open_loginLogic) Open_login(req *types.OpenLoginRequest) (resp *types.LoginResponse, err error) {

	switch req.Flag {
	case "qq":
		info, err := open_login.NewQQLogin(req.Code, open_login.QQConfig{
			AppID:    l.svcCtx.Config.QQ.AppID,
			AppKey:   l.svcCtx.Config.QQ.AppKey,
			Redirect: l.svcCtx.Config.QQ.Redirect,
		})
		if err != nil {
			logx.Error(err)
			return nil, errors.New("qq登录失败")
		}
		log.Println(info)
		var user auth_models.UserModel
		err = l.svcCtx.DB.Take(&user, "open_id=?", info.OpenID).Error
		if err != nil {
			// todo 注册逻辑
			log.Println("注册服务")
		}
		//todo 登录逻辑
		//jwts.GenToken()
	}

	return
}
