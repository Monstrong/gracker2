package httpapi

import "net/http"

// Health godoc
// @Summary Health check
// @Description Checks whether the API Gateway process is running
// @Tags Health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /health [get]
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // необязательно, но пусть будет
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}