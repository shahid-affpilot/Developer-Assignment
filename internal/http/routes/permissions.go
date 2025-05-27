package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/shahid-affpilot/affpilot-auth-service/internal/http/handlers/permission"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
)

func RegisterPermissionRoutes(router *mux.Router) {
	permissions := router.PathPrefix("/api/v1/permission").Subrouter()

	permissions.Handle("",
		middleware.AuthMiddleware(
			middleware.RequireRole("admin", "system_admin")(
				http.HandlerFunc(handlers.GetPermissionList),
			),
		),
	).Methods("GET")

	permissions.Handle("/{permission_id}",
		middleware.AuthMiddleware(
			middleware.RequireRole("admin", "system_admin")(
				http.HandlerFunc(handlers.PermissionDetails),
			),
		),
	).Methods("GET")
}
