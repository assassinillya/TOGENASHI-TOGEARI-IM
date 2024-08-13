package svc

import (
	"gorm.io/gorm"
	"im_server/core"
	"im_server/im_file/file_rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
	}
}
