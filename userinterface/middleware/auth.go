package middleware

import (
	"context"
	"net/http"
	"strings"
	"user-service/application/service"
	"user-service/domain/entities"
)

const AuthContextKey string = "authKey"

func Auth(authService service.AuthService) func (next http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, prefix)
			if token == "" {
				http.Error(w, "Token is empty", http.StatusUnauthorized)
				return
			}

			user, err := authService.Get(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AuthContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) entities.User {
	user, _ := ctx.Value(AuthContextKey).(entities.User)
	return user
}
