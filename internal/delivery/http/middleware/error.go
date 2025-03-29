package middleware

import (
	"encoding/json"
	"net/http"

	"money_tracker/internal/delivery/http/dto"
	"money_tracker/internal/pkg/logger"
)

func RecoveryMiddleware(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.ErrorContext(r.Context(), "panic recovered", map[string]interface{}{
						"error":  err,
						"method": r.Method,
						"path":   r.URL.Path,
					})

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(dto.ErrorResponse{
						Message: "An unexpected error occurred",
					})
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
