package middleware

import (
	"context"
	"net/http"
	"strings"

	"workout-trainer/internal/api"
	"workout-trainer/internal/auth"

	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Authenticate(authenticator *auth.JWTAuthenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				api.ErrorJSON(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				api.ErrorJSON(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			claims, err := authenticator.ValidateToken(parts[1])
			if err != nil {
				api.ErrorJSON(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}
