package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"im_server/common/etcd"
	"im_server/common/middleware"
	"im_server/im_logs/logs_api/internal/mqs"

	"im_server/im_logs/logs_api/internal/config"
	"im_server/im_logs/logs_api/internal/handler"
	"im_server/im_logs/logs_api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/logs.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	// 设置全局中间件
	server.Use(middleware.LogMiddleware)

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range mqs.Consumers(c, context.Background(), ctx) {
		serviceGroup.Add(mq)
	}

	go serviceGroup.Start()
	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
