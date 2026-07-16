package server

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

	pb "github.com/monstrong/gracker2/proto/gen/go/auth/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/auth-service/internal/config"
	"github.com/monstrong/gracker2/auth-service/internal/interceptors"
	"github.com/monstrong/gracker2/auth-service/pkg/db"
	"github.com/monstrong/gracker2/auth-service/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		panic(fmt.Sprintf("error creating db pool and doing Ping: %w", err)) // panic вместо log.Fatal чтобы defer выполнился.
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
	pb.RegisterAuthServiceServer(grpcServer, handler)

	httpServer := infraServerInit(pool, cfg, l)
	g, errgr_ctx := errgroup.WithContext(ctx)

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
        case <-errgr_ctx.Done():
            l.Info(ctx, "errgroup context done, forcing shutdown")
        }

		timeout, cancel := context.WithTimeout(ctx, 15 * time.Second)
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
		go func ()  {
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
			l.Info(ctx, "auth-service stopped gracefully")
		} else {
			l.Error(ctx, "auth-service stopped with error", logger.Error(err))
		}
	} else {
		l.Info(ctx, "auth-service stopped gracefully")
	}
}


func infraServerInit(pool *pgxpool.Pool, cfg *config.Config, l logger.Logger) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("ERROR"))
			l.Error(ctx, "error on connection to db", logger.Error(err))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}
	})

	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Infra),
		Handler: mux,
	}
}