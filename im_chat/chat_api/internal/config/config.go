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
	UserRpc zrpc.RpcClientConf
	FileRpc zrpc.RpcClientConf
	Redis   struct {
		Addr string
		Pwd  string
		DB   int
	}
}
