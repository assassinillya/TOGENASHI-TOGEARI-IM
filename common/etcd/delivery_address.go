package etcd

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"im_server/core"
	"im_server/utils/ips"
	"strings"
)

// DeliveryAddress 上送服务地址
func DeliveryAddress(etcdAddr string, ServiceName string, addr string) {
	list := strings.Split(addr, ":")

	if len(list) != 2 {
		logx.Errorf("ip地址错误%s", addr)
		return
	}

	if list[0] == "0.0.0.0" {
		ip := ips.GetIP()
		addr = strings.ReplaceAll(addr, "0.0.0.0", ip)
	}

	client := core.InitEtcd(etcdAddr)
	_, err := client.Put(context.Background(), ServiceName, addr)
	if err != nil {
		logx.Errorf("服务地址上送失败%s", err.Error())
	}
	logx.Infof("服务地址上送成功%s %s", ServiceName, addr)
}

func GetServiceAddr(etcdAddr string, serviceName string) (addr string) {
	client := core.InitEtcd(etcdAddr)
	res, err := client.Get(context.Background(), serviceName)
	if err == nil && len(res.Kvs) > 0 {
		return string(res.Kvs[0].Value)
	}
	logx.Errorf("获取服务地址失败%s", err.Error())
	return ""
}
