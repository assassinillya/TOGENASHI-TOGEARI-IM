package logic

import (
	"context"
	"errors"
	"fmt"
	"im_server/im_auth/auth_api/internal/svc"
	"im_server/im_auth/auth_api/internal/types"
	"im_server/im_auth/auth_models"
	"im_server/utils/jwts"
	"im_server/utils/pwd"

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

	var user auth_models.UserModel
	err = l.svcCtx.DB.Take(&user, "id = ?", req.UserName).Error
	l.svcCtx.ActionLogs.Info("用户登录操作")
	l.svcCtx.ActionLogs.SetItem("userName", req.UserName)

	defer l.svcCtx.ActionLogs.Save(l.ctx)
	if err != nil {
		l.svcCtx.ActionLogs.Err(req.UserName + "用户名不存在")
		err = errors.New("用户名不存在")
		return
	}

	if pwd.CheckPwd(user.Pwd, req.Password) == false {
		//logx.Error(err)
		l.svcCtx.ActionLogs.SetItem("password", req.Password)
		l.svcCtx.ActionLogs.Err("用户" + req.UserName + " 密码错误")
		err = errors.New("用户名或密码错误")
		return
	}

	token, err := jwts.GenToken(jwts.JwtPayLoad{
		UserID:   user.ID,
		Nickname: user.Nickname,
		Role:     user.Role,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Error(err)
		l.svcCtx.ActionLogs.SetItem("error", err.Error())
		l.svcCtx.ActionLogs.Err("服务内部错误")
		err = errors.New("服务内部错误")
		return
	}

	ctx := context.WithValue(l.ctx, "userID", fmt.Sprintf("%d", user.ID))

	//er := log_stash.NewActionPusher(ctx, l.svcCtx.KqPusherClient, l.svcCtx.Config.Name)
	//er.Info(fmt.Sprintf("%s 用户登录成功", user.Nickname), "")
	//er.Save()
	l.svcCtx.ActionLogs.Info("用户" + user.Nickname + " 登录成功")
	l.svcCtx.ActionLogs.SetCtx(ctx)

	return &types.LoginResponse{Token: token}, nil
}
