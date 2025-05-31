package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/shahid-affpilot/affpilot-auth-service/internal/http/handlers/auth"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
)

func RegisterAuthRoutes(router *mux.Router) {
	auth := router.PathPrefix("/api/v1/auth").Subrouter()
	
	auth.HandleFunc("/register", handlers.Register).Methods("POST")
	auth.HandleFunc("/login", handlers.Login).Methods("POST")

	auth.Handle("/logout",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.Logout),
		),
	).Methods("POST")

	auth.HandleFunc("/verify", handlers.VerifyEmail).Methods("GET")

	auth.Handle("/resend-verification",
		middleware.AuthMiddleware(
			http.HandlerFunc(handlers.ResendVerification),
		),
	).Methods("POST")

	auth.HandleFunc("/password-reset", handlers.InitiatePasswordReset).Methods("POST")
	auth.HandleFunc("/verify/password-reset", handlers.ConfirmPasswordReset).Methods("POST")
}
