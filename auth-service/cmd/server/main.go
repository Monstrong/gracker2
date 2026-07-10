package server

import (
	"context"
	"log"
	"os"

	"github.com/monstrong/auth-service/internal/config"
	"github.com/monstrong/auth-service/pkg/db"
	"github.com/monstrong/auth-service/pkg/logger"
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
		os.Exit(1)
	}
	defer pool.Close()
	l.Info(ctx, "db pool init complete")

	

}