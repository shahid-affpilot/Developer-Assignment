package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/config"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/handlers"
)

var cnf *config.Config

func init() {
	var err error
	cnf, err = config.LoadConfig() // get .env variables
	if err != nil {
		log.Println("Config func does not working well")
	}

	// connect to database
	database.ConnDB(cnf.Database)

	// add system admin if does not exist
	database.InitAdminUser(cnf.Admin)
}

func main() {
	fmt.Println("AffPilot Auth Service starting...")
	log.Println("Server initialized")

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", handlers.Register).Methods("POST")
	auth.HandleFunc("/login", handlers.Login).Methods("POST")
	auth.HandleFunc("/logout", handlers.Logout).Methods("POST")
	auth.HandleFunc("/verify", handlers.VerifyEmail).Methods("GET")
	auth.HandleFunc("/resend-verification", handlers.ResendVerification).Methods("POST")
	auth.HandleFunc("/password-reset", handlers.InitiatePasswordReset).Methods("POST")
	auth.HandleFunc("/verify/password-reset", handlers.ConfirmPasswordReset).Methods("GET")

	// Roles routes
	roles := api.PathPrefix("/roles").Subrouter()
	roles.HandleFunc("", handlers.ListRoles).Methods("GET")
	roles.HandleFunc("/{role_id}", handlers.GetRole).Methods("GET")
	roles.HandleFunc("", handlers.CreateRole).Methods("POST")
	roles.HandleFunc("/{role_id}", handlers.DeleteRole).Methods("DELETE")

	// Permission routes
	permissions := api.PathPrefix("/permission").Subrouter()
	permissions.HandleFunc("", handlers.GetPermissionList).Methods("GET")
	permissions.HandleFunc("/{permission_id}", handlers.PermissionDetails).Methods("GET")

	// Self-user routes
	self := api.PathPrefix("/me").Subrouter()
	self.HandleFunc("", handlers.UserInfo).Methods("GET")
	self.HandleFunc("/permissions", handlers.GetUserPermissions).Methods("GET")

	// users routes
	users := api.PathPrefix("/users").Subrouter()
	users.HandleFunc("", handlers.Users).Methods("GET")
	users.HandleFunc("/{user_id}", handlers.UserDetail).Methods("GET")
	users.HandleFunc("/{user_id}", handlers.UserUpdate).Methods("PUT")
	users.HandleFunc("/{user_id}/request-deletion", handlers.UserDeletion).Methods("POST")
	users.HandleFunc("/{user_id}", handlers.UserDelete).Methods("DELETE")
	users.HandleFunc("/{user_id}/role", handlers.UserRoleChange).Methods("POST")
	users.HandleFunc("/{user_id}/promote/admin", handlers.PromoteToAdmin).Methods("POST")
	users.HandleFunc("/{user_id}/promote/moderator", handlers.PromoteToModerator).Methods("POST")
	users.HandleFunc("/{user_id}/demote", handlers.UserDemote).Methods("POST")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
