package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
)

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// check is there any user logged in or not
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Verification token is required", http.StatusBadRequest)
		return
	}

	// Find user with the given verification token
	var userID uuid.UUID
	var tokenExpiry time.Time
	err := database.DB.QueryRow(`
        SELECT id, token_expiry 
        FROM users 
        WHERE verification_token = $1 
        AND email_verified = false`,
		token,
	).Scan(&userID, &tokenExpiry)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid or expired verification token", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error (cant get userid, token expiry)", http.StatusInternalServerError)
		return
	}

	// Check if token has expired
	if time.Now().After(tokenExpiry) {
		http.Error(w, "Verification token has expired", http.StatusBadRequest)
		return
	}

	// Update user's email_verified status
	_, err = database.DB.Exec(`
        UPDATE users 
        SET email_verified = true,
            verification_token = NULL,
            token_expiry = NULL
        WHERE id = $1`,
		userID,
	)

	if err != nil {
		log.Printf("Error updating user verification status: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return success message
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "verification success",
	})
}
