package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	email "github.com/shahid-affpilot/affpilot-auth-service/internal/services"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func ResendVerification(w http.ResponseWriter, r *http.Request) {
	var req models.ResendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" {
		utils.ErrorResponse(w, http.StatusNoContent, "email is required")
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
		utils.ErrorResponse(w, http.StatusNoContent, "No account is registered with this email")
		return
	}
	if err != nil {
		log.Printf("Database error: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "database error to finde users data")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Account is already verified", nil)

	// Send Email
	verificationURL, _ := email.GenerateVerificationURL(user.ID)
	//str := strings.Split(verificationURL, "=")
	//token := str[1]

	mailText := fmt.Sprintf(
		"Hello %s, here's your new email verification link as requested:\n\n%s\n\nIf you already verified, you can ignore this.\n\nThank you,\nAffpilot AI Team",
		user.Username,
		verificationURL,
	)

	fmt.Println("333")
	err = email.SendVerificationEmail(req.Email, user.Username, mailText)
	if err != nil {
		log.Printf("Error sending verification email: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to send verification email")
		return
	}

	utils.SuccessResponse(w, http.StatusAccepted, "Verification mail has been sent", nil)
}
