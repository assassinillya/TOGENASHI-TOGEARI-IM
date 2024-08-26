package svc

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
	"im_server/core"
	"im_server/im_chat/chat_api/internal/config"
	"im_server/im_file/file_rpc/files"
	"im_server/im_file/file_rpc/types/file_rpc"
	"im_server/im_user/user_rpc/types/user_rpc"
	"im_server/im_user/user_rpc/users"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	UserRpc user_rpc.UsersClient
	FileRpc file_rpc.FilesClient
	Redis   *redis.Client
}

func ClientInfoInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 请求之前
	md := metadata.New(map[string]string{"clientIP": ctx.Value("clientIP").(string), "userID": ctx.Value("userID").(string)})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	err := invoker(ctx, method, req, reply, cc, opts...)
	// 请求之后
	return err
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Pwd, c.Redis.DB)
	return &ServiceContext{
		Config:  c,
		DB:      mysqlDb,
		UserRpc: users.NewUsers(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(ClientInfoInterceptor))),
		FileRpc: files.NewFiles(zrpc.MustNewClient(c.FileRpc)),
		Redis:   client,
	}
}
