package app

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/pkg/config"
	"github.com/monstrong/gracker2/pkg/logger"
	"github.com/monstrong/gracker2/torrent-service/internal/transport/grpc"
)

type App struct {
	cfg *config.Config
	l logger.Logger
	pool *pgxpool.Pool
	grpcServer *grpc.Server
	httpServer *http.Server

}

func New(ctx context.Context) *App {
	
}