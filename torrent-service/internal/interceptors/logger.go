package interceptors

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/torrent-service/pkg/logger"
	"google.golang.org/grpc"
)

const DurationKey = "request-duration"

func LoggingInterseptor(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, 
		req any, 
		info *grpc.UnaryServerInfo, 
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		start := time.Now()
		reqID := uuid.New().String()
		ctx = logger.WithRequestID(ctx, reqID)

		res, err := handler(ctx, req)
		duration := time.Since(start)

		
		if err != nil {
			l.Error(ctx, info.FullMethod,
				logger.Duration(DurationKey, duration),
				logger.Error(err),
			)
		} else {
			l.Info(ctx, info.FullMethod,
				logger.Duration(DurationKey, duration),
			)
		}

		return res, err
	}
}

