package zrpc_interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientInfoInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 请求之前
	md := metadata.New(map[string]string{"clientIP": ctx.Value("clientIP").(string), "userID": ctx.Value("userID").(string)})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	err := invoker(ctx, method, req, reply, cc, opts...)
	// 请求之后
	return err
}
