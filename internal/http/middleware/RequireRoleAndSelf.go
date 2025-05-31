package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func RequireRoleOrSelf(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			userIdFromParam, _ := uuid.Parse(vars["user_id"])

			userId, _ := r.Context().Value(UserIDKey).(string)

			if userId == userIdFromParam.String() {
				next.ServeHTTP(w, r)
				return
			}

			role, ok := r.Context().Value(UserRoleKey).(string)
			if !ok {
				utils.ErrorResponse(w, http.StatusUnauthorized, "no token found")
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
