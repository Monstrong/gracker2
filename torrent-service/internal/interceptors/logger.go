package interceptors

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
	"github.com/monstrong/gracker2/torrent-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const DurationKey = "request-duration"

func LoggingInterceptor(l logger.Logger) grpc.UnaryServerInterceptor {
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

		
		// ErrInvalidData отдается без "обертки"
		if err != nil {
			switch {
			case errors.Is(err, models.ErrNotFound):
				l.Info(ctx, info.FullMethod, logger.Duration(DurationKey, duration), logger.Error(err))
				err = status.Error(codes.NotFound, models.ErrNotFound.Error())
			case errors.Is(err, models.ErrInvalidData):
				l.Info(ctx, info.FullMethod, logger.Duration(DurationKey, duration), logger.Error(err))
				err = status.Error(codes.InvalidArgument, err.Error())
			default:
				l.Error(ctx, info.FullMethod, logger.Duration(DurationKey, duration), logger.Error(err))
				err = status.Error(codes.Internal, models.ErrInternal.Error())
			}
		} else {
			l.Info(ctx, info.FullMethod, logger.Duration(DurationKey, duration),
			)
		}
		return res, err
	}
}

