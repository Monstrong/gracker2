package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/logger"
)

func RequestLogger(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timeStart := time.Now()

			requestID := uuid.NewString()

			traceID := r.Header.Get(string(logger.LoggerTraceIDKey))

			if _, err := uuid.Parse(traceID); err != nil {
				traceID = uuid.NewString()
			}

			ctx := logger.WithRequestID(r.Context(), requestID)
			ctx = logger.WithTraceID(ctx, traceID)

			r = r.WithContext(ctx)

			w.Header().Set(string(logger.LoggerRequestIDKey), requestID)
			w.Header().Set(string(logger.LoggerTraceIDKey), traceID)

			sw := &statusWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(sw, r)

			l.Info(r.Context(), "request processed",
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.Int("status", sw.status),
				logger.Int("bytes", sw.bytes),
				logger.Duration("duration", time.Since(timeStart)),
			)
		})
	}
}


type statusWriter struct {
	http.ResponseWriter //embeding
	status      int
	bytes       int
	wroteHeader bool
}

func (w *statusWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}

	w.status = status
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	n, err := w.ResponseWriter.Write(data)
	w.bytes += n

	return n, err
}

func (w *statusWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

