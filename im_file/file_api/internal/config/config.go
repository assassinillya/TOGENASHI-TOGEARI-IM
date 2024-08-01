package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Etcd      string
	FileSize  float64
	WhiteList []string
}
