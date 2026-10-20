package db

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/pkg/config"
)

const defaultTimeout = 5*time.Second

func NewPool(cfg *config.Postgres) (*pgxpool.Pool, error) {
	u := &url.URL{
		Scheme:   "postgres",
		Host:     net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		User:     url.UserPassword(cfg.User, cfg.Password),
		Path:     cfg.DBName,
	}
	q := u.Query()
	q.Set("sslmode", cfg.SSLMode)
	q.Set("pool_max_conns", strconv.Itoa(cfg.MaxConns))
	u.RawQuery = q.Encode()

	poolCfg, err := pgxpool.ParseConfig(u.String())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("creating pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
