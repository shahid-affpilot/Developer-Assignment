package routes

import (
	"github.com/gorilla/mux"
	handlers "github.com/shahid-affpilot/affpilot-auth-service/internal/http/handlers/user"
)

func RegisterUserRoutes(router *mux.Router) {
	// Self-user routes
	self := router.PathPrefix("/api/v1/me").Subrouter()
	self.HandleFunc("", handlers.UserInfo).Methods("GET")
	self.HandleFunc("/permissions", handlers.GetUserPermissions).Methods("GET")

	// Users management routes
	users := router.PathPrefix("/api/v1/users").Subrouter()
	users.HandleFunc("", handlers.Users).Methods("GET")
	users.HandleFunc("/{user_id}", handlers.UserDetail).Methods("GET")
	users.HandleFunc("/{user_id}", handlers.UserUpdate).Methods("PUT")
	users.HandleFunc("/{user_id}/request-deletion", handlers.UserDeletion).Methods("POST")
	users.HandleFunc("/{user_id}", handlers.UserDelete).Methods("DELETE")
	users.HandleFunc("/{user_id}/role", handlers.UserRoleChange).Methods("POST")
	users.HandleFunc("/{user_id}/promote/admin", handlers.PromoteToAdmin).Methods("POST")
	users.HandleFunc("/{user_id}/promote/moderator", handlers.PromoteToModerator).Methods("POST")
	users.HandleFunc("/{user_id}/demote", handlers.UserDemote).Methods("POST")
}
