package routes

import (
	"github.com/gorilla/mux"
)

// SetupRoutes initializes all routes for the application
func SetupRoutes(router *mux.Router) {
	RegisterAuthRoutes(router)
	RegisterRoleRoutes(router)
	RegisterPermissionRoutes(router)
	RegisterUserRoutes(router)
}
