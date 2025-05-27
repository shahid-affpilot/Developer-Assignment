package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Get user from database
	var user models.User
	err := database.DB.QueryRow(`
        SELECT id, email, password_hash, email_verified, user_type, active
        FROM users
        WHERE email = $1
    `, req.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.EmailVerified, &user.UserType, &user.Active)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !user.Active {
		http.Error(w, "Account is deactivated", http.StatusForbidden)
		return
	}

	// Check if email is verified
	// if !user.EmailVerified {
	// 	http.Error(w, "Please verify your email before logging in", http.StatusUnauthorized)
	// 	return
	// }

	// Verify password
	salt_pass := os.Getenv("PASSWORD_SALT")
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password+salt_pass))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	expiry := os.Getenv("JWT_EXPIRY")
	exp, _ := strconv.Atoi(expiry)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"user_type": user.UserType,
		"exp":       time.Now().Add(time.Duration(exp) * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	
	// Set JWT token in HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   exp * 60 * 60, // exp came from .env
	})

	response := models.LoginResponse{
		UserID:    user.ID.String(),
		Email:     user.Email,
		UserType:  user.UserType,
		ExpiresIn: 24 * 60 * 60, // 24 hours in seconds
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    strconv.Itoa(http.StatusAccepted),
		"message":   "User logged-in success",
		"user_type": string(response.UserType),
		"data":      response,
	})
}
