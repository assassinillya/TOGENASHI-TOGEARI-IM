package logic

import (
	"context"
	"im_server/im_user/user_models"

	"im_server/im_user/user_rpc/internal/svc"
	"im_server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserListInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListInfoLogic {
	return &UserListInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserListInfoLogic) UserListInfo(in *user_rpc.UserListInfoRequest) (*user_rpc.UserListInfoResponse, error) {

	var userList []user_models.UserModel
	l.svcCtx.DB.Find(&userList, in.UserIdList)
	resp := new(user_rpc.UserListInfoResponse)
	resp.UserInfo = make(map[uint32]*user_rpc.UserInfo, 0)
	for _, user := range userList {
		resp.UserInfo[uint32(user.ID)] = &user_rpc.UserInfo{
			NickName: user.Nickname,
			Avatar:   user.Avatar,
		}
	}

	return resp, nil
}
