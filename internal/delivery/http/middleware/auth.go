package middleware

import (
	"net/http"
	"strings"

	httpinternal "money_tracker/internal/delivery/http/internal"
	"money_tracker/internal/domain"
	"money_tracker/internal/pkg/credscontext"
)

type JWTValidator interface {
	ValidateToken(token string) (userID string, role domain.UserRole, err error)
}

func AuthMiddleware(validator JWTValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpinternal.RespondWithError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				httpinternal.RespondWithError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			token := parts[1]
			userID, role, err := validator.ValidateToken(token)
			if err != nil {
				httpinternal.RespondWithError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := credscontext.WithUserID(r.Context(), userID)
			ctx = credscontext.WithUserRole(ctx, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
