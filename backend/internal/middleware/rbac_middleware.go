package middleware

import (
	"leave-management-system/internal/utils"
	"net/http"
)

// RBACMiddleware checks if the user has required roles
func RBACMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user role from context
			role, err := GetUserRoleFromContext(r.Context())
			if err != nil {
				utils.ErrorResponse(w, "Unauthorized", nil, http.StatusUnauthorized)
				return
			}

			// Check if user role is in allowed roles
			hasAccess := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				utils.ErrorResponse(w, "Forbidden: insufficient permissions", nil, http.StatusForbidden)
				return
			}

			// User has required role, proceed
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin middleware allows only admin users
func RequireAdmin(next http.Handler) http.Handler {
	return RBACMiddleware("admin")(next)
}

// RequireAdminOrManager middleware allows admin or manager users
func RequireAdminOrManager(next http.Handler) http.Handler {
	return RBACMiddleware("admin", "manager")(next)
}

// RequireEmployee middleware allows any authenticated user (employee, manager, or admin)
func RequireEmployee(next http.Handler) http.Handler {
	return RBACMiddleware("employee", "manager", "admin")(next)
}
