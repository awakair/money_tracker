package middleware

import (
	"net/http"
	"time"

	"money_tracker/internal/pkg/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += n
	return n, err
}

func LoggerMiddleware(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			if wrapped.statusCode >= 500 {
				logger.ErrorContext(r.Context(), "HTTP request completed with server error",
					"method", r.Method,
					"path", r.URL.Path,
					"status", wrapped.statusCode,
					"duration_ms", duration.Milliseconds(),
					"bytes", wrapped.written,
					"remote_addr", r.RemoteAddr,
				)
			} else {
				logger.InfoContext(r.Context(), "HTTP request completed",
					"method", r.Method,
					"path", r.URL.Path,
					"status", wrapped.statusCode,
					"duration_ms", duration.Milliseconds(),
					"bytes", wrapped.written,
					"remote_addr", r.RemoteAddr,
				)
			}
		})
	}
}
