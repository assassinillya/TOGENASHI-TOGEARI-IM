package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Etcd      string
	FileSize  float64  // 文件大小限制 单位MB
	WhiteList []string // 图片上传的白名单
	BlackList []string // 文件上传的黑名单
	UploadDir string   // 上传文件保存的目录
	UserRpc   zrpc.RpcClientConf
	Mysql     struct {
		DataSource string
	}
}
