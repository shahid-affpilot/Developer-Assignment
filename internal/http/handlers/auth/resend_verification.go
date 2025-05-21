package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/config"
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

	// Check if already verified
	if user.EmailVerified {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Email is already verified",
		})
		return
	}

	// Generate new verification token
	verificationToken := uuid.New().String()
	tokenExpiry := time.Now().Add(5 * time.Minute)

	// Update user with new token
	_, err = database.DB.Exec(`
        UPDATE users 
        SET verification_token = $1,
            token_expiry = $2
        WHERE id = $3`,
		verificationToken,
		tokenExpiry,
		user.ID,
	)

	if err != nil {
		log.Printf("Error updating verification token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send new verification email
	verificationURL := fmt.Sprintf("%s?token=%s", config.EMAIL_VERIFICATION_URL, verificationToken)

	err = email.SendVerificationEmail(req.Email, user.Username, verificationURL)
	if err != nil {
		log.Printf("Error sending verification email: %v", err)
		http.Error(w, "Error sending verification email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "New verification email has been sent",
	})
}
