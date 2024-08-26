package zrpc_interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientInfoInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 请求之前
	var clientIP, userID string
	cl := ctx.Value("clientIP")
	if cl != nil {
		clientIP = cl.(string)
	}
	ui := ctx.Value("userID")
	if ui != nil {
		userID = ui.(string)
	}
	md := metadata.New(map[string]string{"clientIP": clientIP, "userID": userID})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	err := invoker(ctx, method, req, reply, cc, opts...)
	// 请求之后
	return err
}
