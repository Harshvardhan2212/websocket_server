package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"Harshvardhan2212/websocket_server/internal/realtime"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	ClientID contextKey = "clientID"
	Role     contextKey = "role"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		clientID, _ := claims["user_id"].(string)
		role, _ := claims["role"].(string)

		ctx := context.WithValue(r.Context(), ClientID, clientID)
		ctx = context.WithValue(ctx, Role, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contains(slice []realtime.RoleName, target string) bool {
	for _, v := range slice {
		if string(v) == target {
			return true
		}
	}
	return false
}

func RequireRole(roles ...realtime.RoleName) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleVal := r.Context().Value(Role)
			roleStr, ok := roleVal.(string)
			if !ok {
				http.Error(w, "Access Denied", http.StatusForbidden)
				return
			}

			if !contains(roles, roleStr) {
				http.Error(w, "Permission Denied", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
