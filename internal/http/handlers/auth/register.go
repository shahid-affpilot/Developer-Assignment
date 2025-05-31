package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	email "github.com/shahid-affpilot/affpilot-auth-service/internal/services"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		utils.ErrorResponse(w, http.StatusNoContent, "Missing required fields")
		return
	}

	if len(req.Password) < 8 || len(req.Password) > 20 {
		utils.ErrorResponse(w, http.StatusBadRequest, "Password should be in between 8-20 character")
		return
	}

	// checking if username already exists
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", req.Username).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking username: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Database error in checking username")
		return
	}
	if exists {
		utils.ErrorResponse(w, http.StatusConflict, "Username already exists")
		return
	}

	// checking if email already exists
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", req.Email).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking email: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Database error checking email")
		return
	}
	if exists {
		utils.ErrorResponse(w, http.StatusConflict, "Email already registered!")
		return
	}

	// hashing password
	salt := os.Getenv("PASSWORD_SALT")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		utils.ErrorResponse(w, http.StatusForbidden, "Error to hashing pasword")
		return
	}

	// insert into 'users' table
	var user models.User
	err = database.DB.QueryRow(`
		INSERT INTO users (username, email, password_hash, first_name, last_name, email_verified, user_type)
		VALUES ($1, $2, $3, $4, $5, false, 'user')
		RETURNING id, username, email, created_at
	`, req.Username, req.Email, string(hashedPassword), req.FirstName, req.LastName).
		Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		log.Printf("Error inserting user: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Database error in inserting user")
		return
	}

	// insert into 'user-roles' table -- (user & role relationship)
	var defaultRoleID uuid.UUID
	_ = database.DB.QueryRow(`SELECT id FROM roles WHERE name = 'user' LIMIT 1`).Scan(&defaultRoleID)

	_, _ = database.DB.Exec(`INSERT INTO user_roles (user_id, role_id, assigned_by) VALUES ($1, $2, $1)`,
		user.ID, defaultRoleID)

	// Send Email
	verificationURL, _ := email.GenerateVerificationURL(user.ID)

	mailText := fmt.Sprintf(
		"Hello %s,\nYou're registered successfully! Please verify your email using the link below:\n\n%s\n\nThank you,\nAffpilot AI Team",
		user.Username,
		verificationURL,
	)

	err = email.SendVerificationEmail(user.Email, user.Username, mailText)
	if err != nil {
		log.Printf("Error sending verification email: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error sending email")
		return
	}

	// return response
	resp := models.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	utils.SuccessResponse(w, http.StatusCreated, "User registration successful", resp)
}
