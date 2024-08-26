package zrpc_interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ServerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	clientIP := metadata.ValueFromIncomingContext(ctx, "clientIP")
	userID := metadata.ValueFromIncomingContext(ctx, "userID")
	if len(clientIP) > 0 {
		ctx = context.WithValue(ctx, "clientIP", clientIP[0])
	}
	if len(userID) > 0 {
		ctx = context.WithValue(ctx, "userID", userID[0])
	}
	return handler(ctx, req)
}
