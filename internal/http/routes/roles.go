package routes

import (
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/handlers/role"
)

func RegisterRoleRoutes(router *mux.Router) {
	roles := router.PathPrefix("/api/v1/roles").Subrouter()
	roles.HandleFunc("", handlers.ListRoles).Methods("GET")
	roles.HandleFunc("/{role_id}", handlers.GetRole).Methods("GET")
	roles.HandleFunc("", handlers.CreateRole).Methods("POST")
	roles.HandleFunc("/{role_id}", handlers.DeleteRole).Methods("DELETE")
}
