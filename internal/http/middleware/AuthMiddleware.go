package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserRoleKey contextKey = "user_role"
	UserIDKey   contextKey = "user_id"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized - No token", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract required claims
		userRole, ok1 := claims["user_type"].(string)
		userID, ok2 := claims["user_id"].(string)

		if !ok1 || !ok2 {
			http.Error(w, "Unauthorized - Invalid token data", http.StatusUnauthorized)
			return
		}

		// Add to context
		ctx := context.WithValue(r.Context(), UserRoleKey, userRole)
		ctx = context.WithValue(ctx, UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
