package middleware

import (
	"net/http"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(UserRoleKey).(string)

			if !ok {
				http.Error(w, "Unauthorized - No role found", http.StatusUnauthorized)
				utils.ErrorResponse(w, http.StatusUnauthorized, "Unathorized - Ro role found")
				return
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			utils.ErrorResponse(w, http.StatusUnauthorized, "Permission denied")
		})
	}
}
