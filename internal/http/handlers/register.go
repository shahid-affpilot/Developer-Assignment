package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// Basic validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Check if username already exists
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

	// Check if email already exists
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

	// Hash password
	salt_pass := os.Getenv("PASSWORD_SALT")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password+salt_pass), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generate verification token
	verificationToken := uuid.New().String()
	tokenExpiry := time.Now().Add(5 * time.Minute)

	// Create user (insert into database)
	var user models.User
	err = database.DB.QueryRow(`
        INSERT INTO users (
            username, email, password_hash, first_name, last_name,
            email_verified, user_type, verification_token, token_expiry
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, username, email, created_at
    `,
		req.Username,
		req.Email,
		string(hashedPassword),
		req.FirstName,
		req.LastName,
		false,  // email_verified
		"user", // user_type
		verificationToken,
		tokenExpiry,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var defaultRoleID uuid.UUID
	_ = database.DB.QueryRow(`
		SELECT id FROM roles WHERE name = 'user' LIMIT 1
	`).Scan(&defaultRoleID)

	var userRole models.UserRole
	userRole.UserID = user.ID
	userRole.AssignedBy = user.ID
	database.DB.QueryRow(`
		INSERT INTO user_roles (
			user_id, role_id, assigned_by
		)
		VALUES ($1, $2, $1)
	`,
		userRole.UserID, defaultRoleID,
	)

	verificationURL := fmt.Sprintf("%s?token=%s",
		os.Getenv("EMAIL_VERIFICATION_URL"), verificationToken)

	err = email.SendVerificationEmail(user.Email, user.Username, verificationURL)
	if err != nil {
		log.Printf("Error sending verification email: %v", err)
	}

	response := models.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Message:  "Registration successful. Please check your email to verify your account.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
