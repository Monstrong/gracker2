package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/monstrong/gracker2/pkg/logger"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Recoverer(l logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil{
					if p == http.ErrAbortHandler {
						panic(p) 
	// Это специальная сигнальная переменная из стандартной библиотеки net/http
	// Она не является настоящей ошибкой
	// Её смысл — сказать серверу: «прерви обработку этого запроса, не логируй это как проблему». 
					}
					l.Error(r.Context(), "panic recovered",
						logger.String("method", r.Method),
						logger.String("path", r.URL.Path),
						logger.String("panic", fmt.Sprint(p)),
						logger.String("stack", string(debug.Stack())),
					)
					span := trace.SpanFromContext(r.Context())
					span.RecordError(fmt.Errorf("panic: %v", p))
					span.SetStatus(codes.Error, "panic")

					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}