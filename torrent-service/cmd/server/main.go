package main

import (
	"context"
	"log"
	"os"

	"github.com/monstrong/gracker2/torrent-service/internal/config"
	"github.com/monstrong/gracker2/torrent-service/internal/repository"
	"github.com/monstrong/gracker2/torrent-service/internal/service"
	"github.com/monstrong/gracker2/torrent-service/internal/transport/grpc"
	"github.com/monstrong/gracker2/torrent-service/pkg/db"
	"github.com/monstrong/gracker2/torrent-service/pkg/logger"
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
	handler := grpc.NewServer(service)

	
}
