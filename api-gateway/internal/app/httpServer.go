package app

import (
	"fmt"
	"net/http"

	"github.com/monstrong/gracker2/api-gateway/internal/gatewayconfig"
)

func newHttpServer(cfg *gatewayconfig.Server, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: cfg.Timeout,
		WriteTimeout:      cfg.Timeout,
		ReadTimeout:       cfg.Timeout,
		IdleTimeout:       cfg.Timeout * 6,
	}
}
