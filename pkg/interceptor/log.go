package interceptor

import (
	"context"
	"google.golang.org/grpc"
)

func unaryLogInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) {

}
