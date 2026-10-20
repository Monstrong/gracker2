package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"


	"github.com/monstrong/gracker2/pkg/config"
	"github.com/monstrong/gracker2/pkg/db"
	"github.com/monstrong/gracker2/pkg/interceptors"
	"github.com/monstrong/gracker2/pkg/logger"
	pb "github.com/monstrong/gracker2/proto/gen/go/torrent/v1" // <-- при смене версии изменить
	"github.com/monstrong/gracker2/torrent-service/internal/app"
	"github.com/monstrong/gracker2/torrent-service/internal/repository"
	"github.com/monstrong/gracker2/torrent-service/internal/service"
	transport "github.com/monstrong/gracker2/torrent-service/internal/transport/grpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	l, err := logger.New(&cfg.Logger)
	if err != nil {
		log.Fatal(err)
	}

	l.Info(ctx, "config & logger init complete")
	defer func() {
		l.Sync()
	}()

	// makes pool and does Ping
	pool, err := db.NewPool(&cfg.Postgres)
	if err != nil {
		l.Error(ctx, "creating db pool and doing Ping", logger.Error(err))
		panic(fmt.Sprintf("failed to create db pool: %v", err)) // panic вместо log.Fatal чтобы defer выполнился.
	}
	defer pool.Close()
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

	httpServer := app.InfraServerInit(pool, cfg, l)
	g, errgrCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Grpc.Port))
		if err != nil {
			l.Error(ctx, "grpc server starting", logger.Error(err))
			return err
		}
		l.Info(ctx, "grpc server starting", logger.Int("port", cfg.Grpc.Port))
		return grpcServer.Serve(lis)
	})
	g.Go(func() error {
		l.Info(ctx, "http infra server starting")
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-shutdown:
			l.Info(ctx, "signal received, starting graceful shutdown")
		case <-errgrCtx.Done():
			l.Info(ctx, "errgroup context done, forcing shutdown")
		}

		timeout, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
		GrpcStopped := make(chan struct{})
		HttpStopped := make(chan struct{})

		AllStopped := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			l.Info(ctx, "grpc server stopped gracefully")
			close(GrpcStopped)
		}()
		go func() {
			httpServer.Shutdown(context.Background())
			l.Info(ctx, "http server stopped gracefully")
			close(HttpStopped)
		}()
		go func() {
			<-GrpcStopped
			<-HttpStopped

			close(AllStopped)
		}()

		select {
		case <-AllStopped:

		case <-timeout.Done():
			l.Error(ctx, "shut down by timeout")
			grpcServer.Stop()
			httpServer.Close()
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		if err == http.ErrServerClosed {
			l.Info(ctx, "torrent-service stopped gracefully")
		} else {
			l.Error(ctx, "torrent-service stopped with error", logger.Error(err))
		}
	} else {
		l.Info(ctx, "torrent-service stopped gracefully")
	}
}

