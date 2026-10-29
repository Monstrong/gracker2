package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/pkg/config"
	"github.com/monstrong/gracker2/pkg/db"
	"github.com/monstrong/gracker2/pkg/interceptors"
	"github.com/monstrong/gracker2/pkg/logger"
	pb "github.com/monstrong/gracker2/proto/gen/go/torrent/v1" // <-- при смене версии изменить
	"github.com/monstrong/gracker2/torrent-service/internal/repository"
	"github.com/monstrong/gracker2/torrent-service/internal/service"
	transport "github.com/monstrong/gracker2/torrent-service/internal/transport/grpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	l          logger.Logger
	pool       *pgxpool.Pool
	grpcServer *grpc.Server
	httpServer *http.Server
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	l, err := logger.New(&cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	l.Info(ctx, "config & logger init complete")

	// создаёт пул и делает Ping
	pool, err := db.NewPool(&cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}
	l.Info(ctx, "db pool init complete")

	repo := repository.NewPostgresRepository(pool)
	svc := service.NewTorrentService(repo)
	handler := transport.NewServer(svc)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryInterceptor(l),
			interceptors.LoggingInterceptor(l),
			interceptors.TimeoutInterceptor(&cfg.Grpc),
		))
	pb.RegisterTorrentServiceServer(grpcServer, handler)

	httpServer := infraServerInit(pool, &cfg.Infra, l)

	return &App{
		cfg:        cfg,
		l:          l,
		pool:       pool,
		grpcServer: grpcServer,
		httpServer: httpServer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	logCtx := context.WithoutCancel(ctx)

	g, gctx := errgroup.WithContext(ctx)

	// gRPC-сервер
	g.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.Grpc.Port))
		if err != nil {
			return fmt.Errorf("grpc listen: %w", err)
		}
		a.l.Info(logCtx, "grpc server starting", logger.Int("port", a.cfg.Grpc.Port))

		err = a.grpcServer.Serve(lis)
		if errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}
		return err
	})

	// http-сервер (health, ready, metrics)
	g.Go(func() error {
		a.l.Info(logCtx, "http infra server starting")

		err := a.httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil // после Shutdown это не ошибка
		}
		return err
	})

	// Сторож: спит, пока не придёт сигнал или не упадёт один из серверов
	g.Go(func() error {
		<-gctx.Done()
		return a.shutdown(logCtx)
	})

	if err := g.Wait(); err != nil {
		a.l.Error(logCtx, "torrent-service stopped with error", logger.Error(err))
		return err
	}
	a.l.Info(logCtx, "torrent-service stopped gracefully")
	return nil
}

// shutdown останавливает оба сервера одновременно
func (a *App) shutdown(ctx context.Context) error {
	a.l.Info(ctx, "shutting down")

	timeout, cancel := context.WithTimeout(context.Background(), a.cfg.GracefulShutdownTimeout)
	defer cancel()

	grpcStopped := make(chan struct{})
	httpStopped := make(chan struct{})
	allStopped := make(chan struct{})
	

	go func() {
		defer close(grpcStopped)
		a.grpcServer.GracefulStop()

	}()
	var httpErr error
	go func() {
		defer close(httpStopped)
		httpErr = a.httpServer.Shutdown(timeout)
	}()
	go func() {
		<-grpcStopped
		<-httpStopped
		close(allStopped)
	}()

	select {
	case <-allStopped:
		if httpErr != nil {
			a.httpServer.Close()
			return fmt.Errorf("http shutdown: %w", httpErr)
		}
		a.l.Info(ctx, "grpc and http servers stopped gracefully")
		return nil
	case <-timeout.Done():
		a.l.Error(ctx, "shutdown timed out, forcing stop")
		a.grpcServer.Stop()
		a.httpServer.Close()
		return errors.New("shutdown timeout: servers were stopped forcibly")
	}
}

func (a *App) Close() error {
	a.pool.Close()

	a.l.Sync()

	return nil
}