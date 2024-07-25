package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Etcd  string
	Mysql struct {
		DataSource string
	}
	Auth struct {
		AccessSecret string
		AccessExpire int
	}
	Redis struct {
		Addr string
		Pwd  string
		DB   int
	}
	OpenLoginList []struct {
		Name string
		Icon string
		Href string
	}
	QQ struct {
		AppID    string
		AppKey   string
		Redirect string
	}
	UserRpc zrpc.RpcClientConf
}
