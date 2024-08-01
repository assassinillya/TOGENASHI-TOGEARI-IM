package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Etcd      string
	FileSize  float64  // 文件大小限制 单位MB
	WhiteList []string // 图片上传的白名单
	UploadDir string   // 上传文件保存的目录
}
