package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()

			rw := newResponseWriter(w)

			next.ServeHTTP(rw, r)

			logger.Info(
				"HTTP request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status_code", rw.statusCode),
				slog.Duration("duration", time.Since(startedAt)),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("request_id", GetRequestID(r.Context())),
			)
		})
	}
}