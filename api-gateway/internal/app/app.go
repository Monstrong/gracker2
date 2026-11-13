package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/monstrong/gracker2/api-gateway/internal/gatewayconfig"
	"github.com/monstrong/gracker2/pkg/logger"
)

// TODO: при завершении через сигнал Run() не ожидает serveErr после Shutdown. 
// мне гпт сказал что надо это как-то исправить

type App struct {
	cfg        *gatewayconfig.Config
	httpServer *http.Server
	l          logger.Logger
}

func New(ctx context.Context) (*App, error) {
	cfg, err := gatewayconfig.Load()
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	l, err := logger.New(&cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	l.Info(ctx, "config & logger init complete")

	server := newHttpServer(&cfg.Server)

	return &App{
		cfg:        cfg,
		httpServer: server,
		l:          l,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	logCtx := context.WithoutCancel(ctx) // контекст для логгера (останется после отмены основного через gracefulShutdown)

	serveErr := make(chan error, 1)
	go func() {
		a.l.Info(logCtx, "http infra server starting")

		err := a.httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil // штатное завершение
		}
		serveErr <- err
	}()

	select {
	case err := <-serveErr:
		if err != nil {
			a.l.Error(logCtx, "http server failed", logger.Error(err))
			return fmt.Errorf("http server: %w", err)
		}
		a.l.Info(logCtx, "http server stopped")
		return nil
	case <-ctx.Done():
		a.l.Info(logCtx, "shutdown signal received")
	}

	if err := a.shutdown(logCtx); err != nil {
		a.l.Error(logCtx, "graceful shutdown failed", logger.Error(err))
		return err
	}

	a.l.Info(logCtx, "api-gateway stopped gracefully")
    return nil
}

// graceful shutdown
func (a *App) shutdown(ctx context.Context) error {
    a.l.Info(ctx, "shutting down")

    shutCtx, cancel := context.WithTimeout(context.Background(), a.cfg.GracefulShutdownTimeout)
    defer cancel()

    if err := a.httpServer.Shutdown(shutCtx); err != nil {
        _ = a.httpServer.Close()
        return fmt.Errorf("http server shutdown: %w", err)
    }

    return nil
}
func (a *App) Close() error {
	a.l.Sync()

	return nil
}
