package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware creates an auth middleware
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for now - return user from token or allow public access
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// For development, use a default user
				// In production, this would return 401
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Parse and validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Get user ID from claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "invalid user_id in token", http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				http.Error(w, "invalid user_id format", http.StatusUnauthorized)
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
