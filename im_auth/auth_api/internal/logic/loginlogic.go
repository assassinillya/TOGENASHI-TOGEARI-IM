package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"im_server/im_auth/auth_models"
	"im_server/utils/jwts"
	"im_server/utils/pwd"
	"strconv"

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

	var user auth_models.UserModel
	err = l.svcCtx.DB.Take(&user, "id = ?", req.UserName).Error
	if err != nil {
		err = errors.New("用户名或密码错误")
		return
	}

	if pwd.CheckPwd(user.Pwd, req.Password) == false {
		//logx.Error(err)
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
		err = errors.New("服务内部错误")
		return
	}

	ctx := context.WithValue(l.ctx, "userID", fmt.Sprintf("%d", user.ID))

	type Request struct {
		LogType int8   `json:"log_type"` // 日志类型 2 操作日志 3 运行日志
		IP      string `json:"ip"`
		UserID  uint   `json:"user_id"`
		Level   string `json:"level"`
		Title   string `json:"title"`
		Content string `json:"content"` // 日志详情
		Service string `json:"service"` // 服务 记录微服务的名称
	}

	userID := ctx.Value("userID").(string)
	userIntID, _ := strconv.Atoi(userID)

	req1 := Request{
		LogType: 2,
		IP:      ctx.Value("clientIP").(string),
		UserID:  uint(userIntID),
		Level:   "info",
		Title:   fmt.Sprintf("%s 用户登录成功", user.Nickname),
		Content: "xxx",
		Service: l.svcCtx.Config.Name,
	}
	byteData, _ := json.Marshal(req1)

	err = l.svcCtx.KqPusherClient.Push(l.ctx, string(byteData))
	if err != nil {
		logx.Error(err)
	}

	return &types.LoginResponse{Token: token}, nil
}
