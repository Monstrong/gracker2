package logger

import (
	"context"
	"fmt"

	"github.com/monstrong/gracker2/torrent-service/internal/config"
	"go.uber.org/zap"
)

type contextKey string

const (
	loggerRequestIDKey contextKey = "x-request-id"
	loggerTraceIDKey   contextKey = "x-trace-id"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
}

type L struct {
	z *zap.Logger
}

func New(cfg *config.Logger) (Logger, error) {

	var logger *zap.Logger
	var err error
	var invalidLevel bool = false

	switch cfg.Level {
	case "dev", "development", "local":
		logger, err = zap.NewDevelopment()
	case "prod", "production":
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewProduction()
		invalidLevel = true
	}

	if err != nil {
		return nil, fmt.Errorf("creating logger:%w", err)
	}

	if invalidLevel {
		logger.Warn("invalid logger level in config, fallback to production")
		return &L{z: logger}, nil
	}

	return &L{z: logger}, nil
}

func WithRequestID (ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, loggerRequestIDKey, id)
}
func WithTraceID (ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, loggerTraceIDKey, id)
}

func (l *L) extractContextID (ctx context.Context) *zap.Logger{
	logger := l.z
	if reqID, ok := ctx.Value(loggerRequestIDKey).(string); ok {
		logger = logger.With(zap.String("x-request-id", reqID))
	}

	if traceID, ok := ctx.Value(loggerTraceIDKey).(string); ok {
		logger = logger.With(zap.String("x-trace-id", traceID))
	}


	return logger
}

func (l *L) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.extractContextID(ctx).Info(msg, fields...)
}

func (l *L) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.extractContextID(ctx).Warn(msg, fields...)
}

func (l *L) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.extractContextID(ctx).Debug(msg, fields...)
}

func (l *L) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.extractContextID(ctx).Error(msg, fields...)
}
