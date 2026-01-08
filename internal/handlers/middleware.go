package handlers

import (
	"context"
	"net/http"
	"strings"

	"frostbyte-api/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "user_id"

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			token, err := authService.ValidateToken(parts[1])
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				http.Error(w, "User ID not found in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, uint(userIDFloat))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminMiddleware(userService *service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := userService.GetUserByID(userID)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			isAdmin := false
			for _, role := range user.Roles {
				if role.Name == "Admin" {
					isAdmin = true
					break
				}
			}

			if !isAdmin {
				http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(userIDKey).(uint)
	return userID, ok
}
