package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	email "github.com/shahid-affpilot/affpilot-auth-service/internal/services"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 || len(req.Password) > 20 {
		http.Error(w, "Password should be between 8 - 20 characters", http.StatusBadRequest)
		return
	}

	// checking if username already exists
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", req.Username).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking username: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	// checking if email already exists
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", req.Email).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking email: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// hashing password
	salt := os.Getenv("PASSWORD_SALT")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = database.DB.QueryRow(`
		INSERT INTO users (username, email, password_hash, first_name, last_name, email_verified, user_type)
		VALUES ($1, $2, $3, $4, $5, false, 'user')
		RETURNING id, username, email, created_at
	`, req.Username, req.Email, string(hashedPassword), req.FirstName, req.LastName).
		Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		log.Printf("Error inserting user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var defaultRoleID uuid.UUID
	_ = database.DB.QueryRow(`SELECT id FROM roles WHERE name = 'user' LIMIT 1`).Scan(&defaultRoleID)

	_, _ = database.DB.Exec(`INSERT INTO user_roles (user_id, role_id, assigned_by) VALUES ($1, $2, $1)`,
		user.ID, defaultRoleID)

	// Send Email
	verificationURL, _ := email.GenerateVerificationURL(user.ID)

	mailText := fmt.Sprintf(
		"Hello %s, you're registered successfully! Please verify your email using the link below:\n\n%s\n\nThank you,\nAffpilot AI Team",
		user.Username,
		verificationURL,
	)

	err = email.SendVerificationEmail(user.Email, user.Username, mailText)
	if err != nil {
		log.Printf("Error sending verification email: %v", err)
		http.Error(w, "Error sending verification email.", http.StatusBadRequest)
		return
	}

	// return response
	resp := models.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Message:  "Registration successful. Please check your email to verify your account.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  strconv.Itoa(http.StatusCreated),
		"message": "User registration successful",
		"data":    resp,
	})
}
