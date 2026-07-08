package interceptors

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/torrent-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const DurationKey = "request-duration"

func LoggingInterseptor(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, 
		req any, 
		info *grpc.UnaryServerInfo, 
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		reqID := uuid.New().String()
		

		MD, _ := metadata.FromIncomingContext(ctx)
		var traceID string
		vals := MD.Get(string(logger.LoggerTraceIDKey))
		if len(vals) > 0 {
			if _, err := uuid.Parse(vals[0]); err != nil {
				l.Error(ctx, info.FullMethod, logger.Error(err))
				traceID = uuid.New().String()
			} else {
				traceID = vals[0]
			}
		} else {
			traceID = uuid.New().String()
		}


		ctx = logger.WithRequestID(ctx, reqID)
		ctx = logger.WithTraceID(ctx, traceID)
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

