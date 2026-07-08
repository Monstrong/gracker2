package interceptors

import (
	"context"

	"google.golang.org/grpc"
)

func RecoveryInterseptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, 
		req any, 
		info *grpc.UnaryServerInfo, 
		handler grpc.UnaryHandler,
		) (resp any, err error) {

	}
}