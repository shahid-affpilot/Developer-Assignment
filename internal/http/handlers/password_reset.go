package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	email "github.com/shahid-affpilot/affpilot-auth-service/internal/services"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type InitiateResetRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

func InitiatePasswordReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req InitiateResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request format",
		})
		return
	}

	// Generate reset token with encoded password
	randomToken := uuid.New().String()
	encodedPass := base64.URLEncoding.EncodeToString([]byte(req.NewPassword))
	token := randomToken + "." + encodedPass
	expiry := time.Now().Add(15 * time.Minute)

	// Update user with reset token
	result, err := database.DB.Exec(`
        UPDATE users 
        SET verification_token = $1, token_expiry = $2 
        WHERE email = $3`,
		randomToken, expiry, req.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error generating reset token",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "No account exists with this email",
		})
		return
	}

	verificationURL := fmt.Sprintf("%s/password-reset?token=%s",
		os.Getenv("EMAIL_VERIFICATION_URL"), token)

	// Send verification email
	err = email.SendVerificationEmail(req.Email, "Password Reset", verificationURL)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error sending reset email",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset verification email sent",
	})
}

func ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fullToken := r.URL.Query().Get("token")
	if fullToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Reset token is required",
		})
		return
	}

	parts := strings.Split(fullToken, ".")
	if len(parts) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid token format",
		})
		return
	}

	token, encodedPass := parts[0], parts[1]

	// Verify token and get user
	var email string
	err := database.DB.QueryRow(`
        SELECT email 
        FROM users 
        WHERE verification_token = $1 
        AND token_expiry > NOW()`,
		token).Scan(&email)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid or expired reset token",
		})
		return
	}

	// Decode and hash new password
	newPassword, err := base64.URLEncoding.DecodeString(encodedPass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid token",
		})
		return
	}
	fmt.Print(newPassword)

	salt_pass := []byte(os.Getenv("PASSWORD_SALT"))
	pass := append(newPassword, salt_pass...)

	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error processing password",
		})
		return
	}

	_, err = database.DB.Exec(`
        UPDATE users 
        SET password_hash = $1, verification_token = NULL, token_expiry = NULL 
        WHERE email = $2`,
		string(hashedPassword), email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error updating password",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password has been reset successfully",
	})
}
