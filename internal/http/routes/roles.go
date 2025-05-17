package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/shahid-affpilot/affpilot-auth-service/internal/http/handlers/role"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
)

func RegisterRoleRoutes(router *mux.Router) {
	roles := router.PathPrefix("/api/v1/roles").Subrouter()
	roles.Handle("", middleware.AuthMiddleware(middleware.RequireRole("admin", "system_admin")(http.HandlerFunc(handlers.ListRoles)))).Methods("GET")
	roles.Handle("/{role_id}", middleware.AuthMiddleware(middleware.RequireRole("admin", "system_admin")(http.HandlerFunc(handlers.GetRole)))).Methods("GET")
	roles.Handle("", middleware.AuthMiddleware(middleware.RequireRole("admin", "system_admin")(http.HandlerFunc(handlers.CreateRole)))).Methods("POST")
	roles.Handle("/{role_id}", middleware.AuthMiddleware(middleware.RequireRole("admin", "system_admin")(http.HandlerFunc(handlers.DeleteRole)))).Methods("DELETE")
}
