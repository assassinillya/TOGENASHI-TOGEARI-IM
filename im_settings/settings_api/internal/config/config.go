package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Etcd          string
	OpenLoginList []struct {
		Name string
		Icon string
		Href string
	}
}
