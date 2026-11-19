package httpapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/monstrong/gracker2/api-gateway/internal/health"
	"github.com/monstrong/gracker2/api-gateway/internal/transport/httpapi/dto"
)

// Ready godoc
// @Summary Readiness check
// @Description Checks whether the API Gateway and its dependencies are ready to handle requests
// @Tags Health
// @Produce plain
// @Success 200 {string} string "OK"
// @Failure 503 {object} dto.ErrorResponse
// @Router /ready [get]
func Ready(checker *health.Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), checker.GetTimeout())
		defer cancel()

		if err := checker.Check(ctx); err != nil {
			writeError(w, http.StatusServiceUnavailable, "service unavailable")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	_ = writeJSON(w, status, dto.ErrorResponse{
		Error: message,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
