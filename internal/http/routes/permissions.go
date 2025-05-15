package routes

import (
	"github.com/gorilla/mux"
	handlers "github.com/shahid-affpilot/affpilot-auth-service/internal/http/handlers/permission"
)

func RegisterPermissionRoutes(router *mux.Router) {
	permissions := router.PathPrefix("/api/v1/permission").Subrouter()
	permissions.HandleFunc("", handlers.GetPermissionList).Methods("GET")
	permissions.HandleFunc("/{permission_id}", handlers.PermissionDetails).Methods("GET")
}
