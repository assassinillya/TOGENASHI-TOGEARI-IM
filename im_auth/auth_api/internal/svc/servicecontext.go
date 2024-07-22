package svc

import (
	"gorm.io/gorm"
	"im_server/core"
	"im_server/im_auth/auth_api/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	//mysqlDb.AutoMigrate(&auth_models.UserModel{})
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
	}
}
