package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/pkg/config"
	"github.com/monstrong/gracker2/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func infraServerInit(pool *pgxpool.Pool, rclient *redis.Client, cfg *config.Infra, l logger.Logger) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), cfg.Timeout/2)
		defer cancel()

		Rerr := rclient.Ping(ctx).Err()
		Perr := pool.Ping(ctx)
		// err := errors.Join(pool.Ping(ctx), rclient.Ping(ctx).Err()).
		switch {
		case Rerr != nil && Perr != nil:
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("ERROR"))
			l.Error(ctx, "error on connection to db & redis", logger.Error(Perr), logger.Error(Rerr))
		case Rerr != nil:
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("ERROR"))
			l.Error(ctx, "error on connection to redis", logger.Error(Rerr))
		case Perr != nil:
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("ERROR"))
			l.Error(ctx, "error on connection to db", logger.Error(Perr))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}
	})

	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
		ReadHeaderTimeout: cfg.Timeout,
		WriteTimeout:      cfg.Timeout,
		ReadTimeout:       cfg.Timeout,
		IdleTimeout:       cfg.Timeout*6,
	}
}
