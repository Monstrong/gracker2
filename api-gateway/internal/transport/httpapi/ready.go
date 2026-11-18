package httpapi

import (
	"context"
	"net/http"

	"github.com/monstrong/gracker2/api-gateway/internal/health"
)

func Ready(checker *health.Checker) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), checker.GetTimeout())
        defer cancel()

        if err := checker.Check(ctx); err != nil {
            http.Error(
                w,
                "service unavailable",
                http.StatusServiceUnavailable,
            )
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }
}