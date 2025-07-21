package middleware

import (
	"context"
	"me-ai/pkg/jwt"
	"net/http"
	"strings"
)

type contextKey string

const UserEmailKey contextKey = "user_email"

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(header, "Bearer ")
			ok, data := jwt.NewJWT(secret).Parse(token)
			if !ok || data == nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserEmailKey, data.Login)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserEmail(r *http.Request) string {
	if v := r.Context().Value(UserEmailKey); v != nil {
		if email, ok := v.(string); ok {
			return email
		}
	}
	return ""
}
