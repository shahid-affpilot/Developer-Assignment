package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	email "github.com/shahid-affpilot/affpilot-auth-service/internal/services"
)

func ResendVerification(w http.ResponseWriter, r *http.Request) {
	var req models.ResendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Get user details
	var user struct {
		ID            uuid.UUID
		Username      string
		EmailVerified bool
	}

	fmt.Println("111")
	err := database.DB.QueryRow(`
        SELECT id, username, email_verified 
        FROM users 
        WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.Username, &user.EmailVerified)

	if err == sql.ErrNoRows {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("222")
	// Check if already verified
	if user.EmailVerified {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Email is already verified",
		})
		return
	}

	// Send Email
	verificationURL, _ := email.GenerateVerificationURL(user.ID)

	mailText := fmt.Sprintf(
		"Hello %s, here's your new email verification link as requested:\n\n%s\n\nIf you already verified, you can ignore this.\n\nThank you,\nAffpilot AI Team",
		user.Username,
		verificationURL,
	)

	fmt.Println("333")
	err = email.SendVerificationEmail(req.Email, user.Username, mailText)
	if err != nil {
		log.Printf("Error sending verification email: %v", err)
		http.Error(w, "Error sending verification email.", http.StatusBadRequest)
		return
	}

	fmt.Println("444")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "New verification email has been sent",
	})
}
