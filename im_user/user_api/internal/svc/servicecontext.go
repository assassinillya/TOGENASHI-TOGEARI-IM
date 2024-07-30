package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"im_server/im_user/user_api/internal/config"
	"im_server/im_user/user_rpc/types/user_rpc"
	"im_server/im_user/user_rpc/users"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user_rpc.UsersClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
	}
}
