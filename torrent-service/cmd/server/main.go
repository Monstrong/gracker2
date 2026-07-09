package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/monstrong/gracker2/proto/gen/go/torrent/v1" // <-- при смене версии изменить
	"github.com/monstrong/gracker2/torrent-service/internal/config"
	"github.com/monstrong/gracker2/torrent-service/internal/interceptors"
	"github.com/monstrong/gracker2/torrent-service/internal/repository"
	"github.com/monstrong/gracker2/torrent-service/internal/service"
	transport "github.com/monstrong/gracker2/torrent-service/internal/transport/grpc"
	"github.com/monstrong/gracker2/torrent-service/pkg/db"
	"github.com/monstrong/gracker2/torrent-service/pkg/logger"
	"google.golang.org/grpc"
)

func main() {

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("some error loading config")
	}

	l, err := logger.New(&cfg.Logger)
	if err != nil {
		log.Fatal("some error creating logger")
	}
	
	l.Info(ctx, "config & logger init complete")
	defer func() {
        l.Sync()
    }()

	// makes pool and does Ping
	pool, err := db.NewPool(&cfg.Postgres)
	if err != nil {
		l.Error(ctx, "creating db pool and doing Ping", logger.Error(err))
		os.Exit(1)
	}
	defer pool.Close()
	l.Info(ctx, "db pool init complete")

	repo := repository.NewPostgresRepository(pool)
	service := service.NewTorrentService(repo)
	handler := transport.NewServer(service)
	
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.LoggingInterceptor(l), 
			interceptors.RecoveryInterceptor(l),
		))
	pb.RegisterTorrentServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.App.Port))
	if err != nil {
		l.Error(ctx, "server starting", logger.Error(err))
		os.Exit(1)
	}

	go GracefulStop(ctx, l, grpcServer)
	
	l.Info(ctx, "grpc server starting", logger.Int("port", cfg.App.Port))
	if err := grpcServer.Serve(lis); err != nil {
		l.Error(ctx, "grpc server failed", logger.Error(err))
		os.Exit(1)
	}
}

func GracefulStop(ctx context.Context, l logger.Logger, grpcServer *grpc.Server) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdown
		l.Info(ctx, "shutting down gracefully...")
		timeout, cancel := context.WithTimeout(ctx, 15 * time.Second)
		defer cancel()
		stopped := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
			l.Info(ctx, "shut down gracefully")
		case <-timeout.Done():
			l.Error(ctx, "shut down by timeout")
			grpcServer.Stop()
		}
	}()
}