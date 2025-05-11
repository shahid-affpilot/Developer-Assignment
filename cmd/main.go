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
	auth.HandleFunc("/password-reset", handlers.ResetPassword).Methods("POST")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
