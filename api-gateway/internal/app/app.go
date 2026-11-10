package app

import (
	"context"
	"fmt"

	"github.com/monstrong/gracker2/api-gateway/internal/gatewayconfig"
)

type App struct {
	cfg *gatewayconfig.Config
}


func New(ctx context.Context) (*App, error) {
	cfg, err := gatewayconfig.Load()
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	return &App{cfg: cfg}, nil
}