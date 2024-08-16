package svc

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"im_server/core"
	"im_server/im_group/group_rpc/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	DB        *gorm.DB
	RedisConf *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.RedisConf.Addr, c.RedisConf.Pwd, c.RedisConf.DB)
	return &ServiceContext{
		RedisConf: client,
		Config:    c,
		DB:        mysqlDb,
	}
}
