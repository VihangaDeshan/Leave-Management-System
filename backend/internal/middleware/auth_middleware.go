package middleware

import (
	"context"
	"leave-management-system/internal/utils"
	apperrors "leave-management-system/pkg/errors"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"
const UserEmailKey contextKey = "userEmail"

// AuthMiddleware validates JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(w, "Authorization header required", nil, http.StatusUnauthorized)
			return
		}

		// Check for Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(w, "Invalid authorization header format", nil, http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Verify token
		claims, err := utils.VerifyToken(token)
		if err != nil {
			utils.ErrorResponse(w, "Invalid or expired token", nil, http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext retrieves user ID from request context
func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserIDKey).(int)
	if !ok {
		return 0, apperrors.Unauthorized("User ID not found in context")
	}
	return userID, nil
}

// GetUserRoleFromContext retrieves user role from request context
func GetUserRoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value(UserRoleKey).(string)
	if !ok {
		return "", apperrors.Unauthorized("User role not found in context")
	}
	return role, nil
}

// GetUserEmailFromContext retrieves user email from request context
func GetUserEmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value(UserEmailKey).(string)
	if !ok {
		return "", apperrors.Unauthorized("User email not found in context")
	}
	return email, nil
}
