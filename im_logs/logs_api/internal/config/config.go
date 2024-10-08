package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	KqConsumerConf kq.KqConf
	UserRpc        zrpc.RpcClientConf
	Mysql          struct {
		DataSource string
	}
	Etcd         string
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
}
