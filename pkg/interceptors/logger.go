package interceptors

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/errs"
	"github.com/monstrong/gracker2/pkg/logger"
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

		if err != nil {
			var ve *errs.ValidationError
			switch {
			case errors.As(err, &ve):
				l.Info(ctx, info.FullMethod, logger.Duration(DurationKey, duration), logger.Error(err))
			case errors.Is(err, errs.ErrNotFound),
				errors.Is(err, errs.ErrAlreadyExists),
				errors.Is(err, errs.ErrInvalidCredentials),
				errors.Is(err, errs.ErrNameAlreadyExists),
				errors.Is(err, errs.ErrEmailAlreadyExists),
				errors.Is(err, errs.ErrUserPatchEmpty):
				
				l.Info(ctx, info.FullMethod, logger.Duration(DurationKey, duration), logger.Error(err))
			case errors.Is(err, context.DeadlineExceeded),
				errors.Is(err, context.Canceled):
				l.Warn(ctx, info.FullMethod, logger.Duration(DurationKey, duration), logger.Error(err))
			default:
				l.Error(ctx, info.FullMethod, logger.Duration(DurationKey, duration), logger.Error(err))
			}
		} else {
			l.Info(ctx, info.FullMethod, logger.Duration(DurationKey, duration))
		}

		err = toGRPCError(err) // значение меняется на статус код, для return

		return res, err
	}
}

func toGRPCError(err error) error {
	if err == nil {
		return nil
	}

	var ve *errs.ValidationError
	if errors.As(err, &ve) {
		return status.Error(codes.InvalidArgument, ve.Reason)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return status.Error(codes.DeadlineExceeded, err.Error())
	}
	if errors.Is(err, context.Canceled) {
		return status.Error(codes.Canceled, err.Error())
	}
	if errors.Is(err, errs.ErrNotFound) {
		return status.Error(codes.NotFound, errs.ErrNotFound.Error())
	}
	if errors.Is(err, errs.ErrAlreadyExists) {
		return status.Error(codes.AlreadyExists, errs.ErrAlreadyExists.Error())
	}
	if errors.Is(err, errs.ErrUserPatchEmpty) {
		return status.Error(codes.InvalidArgument, errs.ErrUserPatchEmpty.Error())
	}
	if errors.Is(err, errs.ErrInvalidCredentials) { // for auth
		return status.Error(codes.Unauthenticated, errs.ErrInvalidCredentials.Error())
	}
	if errors.Is(err, errs.ErrNameAlreadyExists) { // for auth
		return status.Error(codes.AlreadyExists, errs.ErrNameAlreadyExists.Error())
	}
	if errors.Is(err, errs.ErrEmailAlreadyExists) { // for auth
		return status.Error(codes.AlreadyExists, errs.ErrEmailAlreadyExists.Error())
	}
	return status.Error(codes.Internal, errs.ErrInternal.Error())
}
