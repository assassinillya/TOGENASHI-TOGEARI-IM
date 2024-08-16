package svc

import (
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"im_server/core"
	"im_server/im_group/group_api/internal/config"
	"im_server/im_group/group_rpc/groups"
	"im_server/im_group/group_rpc/types/group_rpc"
	"im_server/im_user/user_rpc/types/user_rpc"
	"im_server/im_user/user_rpc/users"
)

type ServiceContext struct {
	Config   config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	UserRpc  user_rpc.UsersClient
	GroupRpc group_rpc.GroupsClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Pwd, c.Redis.DB)
	return &ServiceContext{
		Config:   c,
		DB:       mysqlDb,
		Redis:    client,
		UserRpc:  users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
		GroupRpc: groups.NewGroups(zrpc.MustNewClient(c.GroupRpc)),
	}
}
