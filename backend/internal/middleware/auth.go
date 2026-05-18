package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/stockapp/backend/internal/service"
)

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for public routes
			publicPaths := []string{"/auth/login", "/auth/register", "/uploads"}
			path := r.URL.Path
			for _, p := range publicPaths {
				if strings.HasPrefix(path, p) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Get token from header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			claims, err := authService.ValidateToken(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Set user ID in context
			sub, _ := (*claims)["sub"].(string)
			role, _ := (*claims)["role"].(string)
			r = setUserID(r, sub)
			r = setUserRole(r, role)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(r *http.Request) string {
	if v := r.Context().Value("userID"); v != nil {
		return v.(string)
	}
	return ""
}

func GetUserRole(r *http.Request) string {
	if v := r.Context().Value("userRole"); v != nil {
		return v.(string)
	}
	return ""
}

func setUserID(r *http.Request, userID string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "userID", userID)
	return r.WithContext(ctx)
}

func setUserRole(r *http.Request, role string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "userRole", role)
	return r.WithContext(ctx)
}