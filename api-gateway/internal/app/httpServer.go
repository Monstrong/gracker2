package app

import (
	"fmt"
	"net/http"

	"github.com/monstrong/gracker2/api-gateway/internal/gatewayconfig"
	"github.com/monstrong/gracker2/api-gateway/internal/transport/httpapi"
)

func newHttpServer(cfg *gatewayconfig.Server) *http.Server {
	r := httpapi.NewRouter

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r(),
		ReadHeaderTimeout: cfg.Timeout,
		WriteTimeout:      cfg.Timeout,
		ReadTimeout:       cfg.Timeout,
		IdleTimeout:       cfg.Timeout*6,
	}
}