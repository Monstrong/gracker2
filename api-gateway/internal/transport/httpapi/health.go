package httpapi

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // необязательно, но пусть будет
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}