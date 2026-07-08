package interceptors

import (
	"context"
	"fmt"

	"github.com/monstrong/gracker2/torrent-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func RecoveryInterseptor(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, 
		req any, 
		info *grpc.UnaryServerInfo, 
		handler grpc.UnaryHandler,
		) (resp any, err error) {

		defer func() {
			if p := recover(); p != nil{
				l.Error(ctx, "panic recovered",
					logger.String("method", info.FullMethod),
					logger.String("panic", fmt.Sprintf("%v", p)),
				)

				// подмена ошибки чтобы пользователь не увидел её текст
				// это - то, для чего здесь возвращаемые параметры проименнованы
				err = status.Error(codes.Internal, "internal server error") 
			}
		}()
		return handler(ctx, req)
	}
}