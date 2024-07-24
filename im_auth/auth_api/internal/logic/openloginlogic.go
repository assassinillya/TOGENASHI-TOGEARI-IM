package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"im_server/im_auth/auth_api/internal/svc"
	"im_server/im_auth/auth_api/internal/types"
	"im_server/im_auth/auth_models"
	"im_server/im_user/user_rpc/types/user_rpc"
	"im_server/utils/jwts"
	"im_server/utils/open_login"
	"log"
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
			res, err := l.svcCtx.UserRpc.UserCreate(context.Background(), &user_rpc.UserCreateRequest{
				NickName: info.Nickname,
				Password: "",
				Role:     2,
				Avatar:   info.Avatar,
				OpenId:   info.OpenID,
			})
			if err != nil {
				logx.Error(err)
				return nil, errors.New("注册失败")
			}
			user.Model.ID = uint(res.UserId)
			user.Role = 2
			user.Nickname = info.Nickname
		}
		//todo 登录逻辑
		token, err := jwts.GenToken(jwts.JwtPayLoad{
			UserID:   user.ID,
			Nickname: user.Nickname,
			Role:     user.Role,
		}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
		if err != nil {
			logx.Error(err)
			err = errors.New("服务内部错误")
			return
		}

		return &types.LoginResponse{Token: token}, nil
	}

	return
}
