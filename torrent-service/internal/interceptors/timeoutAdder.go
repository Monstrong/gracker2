package interceptors

import (
	"context"

	"github.com/monstrong/gracker2/torrent-service/internal/config"
	"google.golang.org/grpc"
)

func TimeoutInterceptor(cfg *config.Grpc) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, 
		req any, 
		info *grpc.UnaryServerInfo, 
		handler grpc.UnaryHandler,
		) (resp any, err error) {

		if _, ok := ctx.Deadline(); ok {
			return handler(ctx, req)
		}

		ctxNew, cancel := context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
		
		return handler(ctxNew, req)
			
	}
}