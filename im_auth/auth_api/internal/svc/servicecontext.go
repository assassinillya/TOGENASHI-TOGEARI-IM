package svc

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"im_server/core"
	"im_server/im_auth/auth_api/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	redisClient := core.InitRedis(c.Redis.Addr, c.Redis.Pwd, c.Redis.DB)

	//mysqlDb.AutoMigrate(&auth_models.UserModel{})
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
		Redis:  redisClient,
	}
}
