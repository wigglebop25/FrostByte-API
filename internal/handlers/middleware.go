package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"frostbyte-api/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

type contextKey string

const (
	userIDKey   contextKey = "user_id"
	usernameKey contextKey = "username"
	rolesKey    contextKey = "roles"
)

// Token Bucket Rate Limiter
type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
}

type visitor struct {
	tokens     float64
	lastUpdate time.Time
}

var limiter = &rateLimiter{
	visitors: make(map[string]*visitor),
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rate := 1.0 // 1 request per second (average)
	burst := 20.0
	now := time.Now()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{tokens: burst - 1, lastUpdate: now}
		return true
	}

	// Refill tokens
	elapsed := now.Sub(v.lastUpdate).Seconds()
	v.tokens += elapsed * rate
	if v.tokens > burst {
		v.tokens = burst
	}
	v.lastUpdate = now

	if v.tokens >= 1 {
		v.tokens--
		return true
	}

	return false
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		if !limiter.allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":       "RateLimitExceeded",
				"message":     "You are sending requests too quickly. Please wait.",
				"retry_after": 5, // Simple static advice
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strict-Transport-Security: max-age=63072000; includeSubDomains; preload (2 years)
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		// XSS Protection (for older browsers)
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Content Security Policy (Basic restricted)
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

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

			// Extract username and roles (if available)
			username, _ := claims["username"].(string)

			var roles []string
			if rolesClaim, ok := claims["roles"].([]interface{}); ok {
				for _, r := range rolesClaim {
					if rStr, ok := r.(string); ok {
						roles = append(roles, rStr)
					}
				}
			}

			ctx := context.WithValue(r.Context(), userIDKey, uint(userIDFloat))
			ctx = context.WithValue(ctx, usernameKey, username)
			ctx = context.WithValue(ctx, rolesKey, roles)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminMiddleware(userService *service.UserService) func(http.Handler) http.Handler {
	return RoleMiddleware(userService, "Admin")
}

func RoleMiddleware(userService *service.UserService, allowedRoles ...string) func(http.Handler) http.Handler {
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

			hasRole := false
			for _, userRole := range user.Roles {
				for _, allowed := range allowedRoles {
					if userRole.Name == allowed {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				http.Error(w, "Unauthorized: Permissions Required", http.StatusForbidden)
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

func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey).(string)
	return username, ok
}

func GetRolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(rolesKey).([]string)
	return roles, ok
}
