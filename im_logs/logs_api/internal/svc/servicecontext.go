package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"im_server/common/zrpc_interceptor"
	"im_server/core"
	"im_server/im_logs/logs_api/internal/config"
	"im_server/im_logs/logs_api/internal/middleware"
	"im_server/im_user/user_rpc/types/user_rpc"
	"im_server/im_user/user_rpc/users"
	"net/http"
)

type ServiceContext struct {
	Config          config.Config
	UserRpc         user_rpc.UsersClient
	DB              *gorm.DB
	AdminMiddleware func(next http.HandlerFunc) http.HandlerFunc
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	return &ServiceContext{
		Config:          c,
		DB:              mysqlDb,
		UserRpc:         users.NewUsers(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		AdminMiddleware: middleware.NewAdminMiddleware().Handle,
	}
}
