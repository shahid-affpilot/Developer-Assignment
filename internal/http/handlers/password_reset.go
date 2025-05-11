package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in by verifying auth cookie
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Please login first",
		})
		return
	}

	// Parse JWT token from cookie
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid session, please login again",
		})
		return
	}

	// Get request body
	var req PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.OldPassword == "" || req.NewPassword == "" {
		http.Error(w, "Old password and new password are required", http.StatusBadRequest)
		return
	}

	if len(req.NewPassword) < 8 {
		http.Error(w, "New password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	// Get user's current password hash from database
	userID := claims["user_id"].(string)
	var currentPasswordHash string
	err = database.DB.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).
		Scan(&currentPasswordHash)

	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Verify old password
	salt_pass := os.Getenv("PASSWORD_SALT")
	err = bcrypt.CompareHashAndPassword(
		[]byte(currentPasswordHash),
		[]byte(req.OldPassword+salt_pass),
	)
	if err != nil {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword+salt_pass),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update password in database
	_, err = database.DB.Exec(
		"UPDATE users SET password_hash = $1 WHERE id = $2",
		string(hashedPassword),
		userID,
	)

	if err != nil {
		log.Printf("Error updating password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password updated successfully",
	})
}
