package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/monstrong/gracker2/torrent-service/internal/config"
	"go.uber.org/zap"
)

type contextKey string

const (
	LoggerRequestIDKey contextKey = "x-request-id"
	LoggerTraceIDKey   contextKey = "x-trace-id"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Debug(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)

}

type L struct {
	z *zap.Logger
}

type Field = zap.Field


func New(cfg *config.Logger) (*L, error) {

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
	return context.WithValue(ctx, LoggerRequestIDKey, id)
}
func WithTraceID (ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, LoggerTraceIDKey, id)
}

func (l *L) extractContextID (ctx context.Context) *zap.Logger{
	logger := l.z
	if reqID, ok := ctx.Value(LoggerRequestIDKey).(string); ok {
		logger = logger.With(zap.String("x-request-id", reqID))
	}

	if traceID, ok := ctx.Value(LoggerTraceIDKey).(string); ok {
		logger = logger.With(zap.String("x-trace-id", traceID))
	}


	return logger
}

func (l *L) Info(ctx context.Context, msg string, fields ...Field) {
	l.extractContextID(ctx).Info(msg, fields...)
}

func (l *L) Warn(ctx context.Context, msg string, fields ...Field) {
	l.extractContextID(ctx).Warn(msg, fields...)
}

func (l *L) Debug(ctx context.Context, msg string, fields ...Field) {
	l.extractContextID(ctx).Debug(msg, fields...)
}

func (l *L) Error(ctx context.Context, msg string, fields ...Field) {
	l.extractContextID(ctx).Error(msg, fields...)
}

func (l *L) Sync() {
	_ = l.z.Sync()
}

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func Error(err error) Field {
	return zap.Error(err)
}