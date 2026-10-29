package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/auth-service/internal/authconfig"
	"github.com/monstrong/gracker2/auth-service/internal/repository"
	"github.com/monstrong/gracker2/auth-service/internal/service"
	transport "github.com/monstrong/gracker2/auth-service/internal/transport/grpc"
	"github.com/monstrong/gracker2/auth-service/pkg/hasher"
	"github.com/monstrong/gracker2/auth-service/pkg/jwtmanager"
	"github.com/monstrong/gracker2/auth-service/pkg/redisclient"
	"github.com/monstrong/gracker2/pkg/db"
	"github.com/monstrong/gracker2/pkg/interceptors"
	"github.com/monstrong/gracker2/pkg/logger"
	pb "github.com/monstrong/gracker2/proto/gen/go/auth/v1" // <-- при смене версии изменить
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *authconfig.Config
	l          logger.Logger
	pool       *pgxpool.Pool
	rclient    *redis.Client
	grpcServer *grpc.Server
	httpServer *http.Server
}

func New(ctx context.Context) (*App, error) {
	cfg, err := authconfig.Load()
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

	rclient, err := redisclient.NewClient(ctx, &cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("create redis client: %w", err)
	}
	l.Info(ctx, "redis client init complete")

	urepo := repository.NewPostgresRepository(pool)
	rfrepo := repository.NewRedisRepository(rclient)
	m := jwtmanager.NewManager(cfg.JWT.SigningKey, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)
	h := hasher.NewSimpleHasher(cfg.Hasher.Cost)
	
	svc := service.New(m, h, urepo, rfrepo)
	handler := transport.NewServer(svc)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryInterceptor(l),
			interceptors.LoggingInterceptor(l),
			interceptors.TimeoutInterceptor(&cfg.Grpc),
		))
	pb.RegisterAuthServiceServer(grpcServer, handler)

	httpServer := infraServerInit(pool, rclient, cfg.Infra.Port, l)

	return &App{
		cfg:        cfg,
		l:          l,
		pool:       pool,
		rclient:    rclient,
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
		a.l.Error(logCtx, "auth-service stopped with error", logger.Error(err))
		return err
	}
	a.l.Info(logCtx, "auth-service stopped gracefully")
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
	a.rclient.Close()

	a.l.Sync()

	return nil
}