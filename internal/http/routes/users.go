package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/shahid-affpilot/affpilot-auth-service/internal/http/handlers/user"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
)

func RegisterUserRoutes(router *mux.Router) {
	// Self-user routes
	self := router.PathPrefix("/api/v1/me").Subrouter()
	self.Handle("", middleware.AuthMiddleware(http.HandlerFunc(handlers.UserInfo))).Methods("GET")
	self.Handle("/permissions", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserPermissions))).Methods("GET")

	// Users management routes
	users := router.PathPrefix("/api/v1/users").Subrouter()
	users.Handle("", middleware.AuthMiddleware(middleware.RequireRole("admin", "system_admin")(http.HandlerFunc(handlers.Users)))).Methods("GET")

	users.Handle("/{user_id}", middleware.AuthMiddleware(middleware.RequireRoleOrSelf("moderator", "admin", "system_admin")(http.HandlerFunc(handlers.UserDetail)))).Methods("GET")

	users.Handle("/{user_id}", middleware.AuthMiddleware(middleware.RequireRoleOrSelf("admin", "system_admin")(http.HandlerFunc(handlers.UserUpdate)))).Methods("PUT")

	users.Handle("/{user_id}/request-deletion", middleware.AuthMiddleware(http.HandlerFunc(handlers.UserDeletion))).Methods("POST")

	users.Handle("/{user_id}", middleware.AuthMiddleware(middleware.RequireRole("moderator", "admin", "system_admin")(http.HandlerFunc(handlers.UserDelete)))).Methods("DELETE")

	users.Handle("/{user_id}/role", middleware.AuthMiddleware(middleware.RequireRole("admin", "system_admin")(http.HandlerFunc(handlers.UserRoleChange)))).Methods("POST")

	users.Handle("/{user_id}/promote/admin", middleware.AuthMiddleware(middleware.RequireRole("system_admin")(http.HandlerFunc(handlers.PromoteToAdmin)))).Methods("POST")

	users.Handle("/{user_id}/promote/moderator", middleware.AuthMiddleware(middleware.RequireRole("admin", "system_admin")(http.HandlerFunc(handlers.PromoteToModerator)))).Methods("POST")

	users.Handle("/{user_id}/demote", middleware.AuthMiddleware(middleware.RequireRole("admin", "system_admin")(http.HandlerFunc(handlers.UserDemote)))).Methods("POST")
}
