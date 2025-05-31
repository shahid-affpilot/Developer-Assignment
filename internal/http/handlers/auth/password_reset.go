package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/config"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	email "github.com/shahid-affpilot/affpilot-auth-service/internal/services"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var Cnf = config.GetConfig()

func InitiatePasswordReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.InitiateResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Generate reset token with encoded password
	cnf := config.GetConfig()
	ttl := time.Duration(cnf.Email.VerificationTTL) * time.Minute
	now := time.Now()

	claims := jwt.MapClaims{
		"email":   req.Email,
		"created": now.Unix(),
		"expiry":  now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(cnf.JWT.Secret))

	verificationURL := fmt.Sprintf("%s/password-reset?token=%s", Cnf.Email.VerificationURL, signedToken)
	mailText := fmt.Sprintf(
		"Hello dear,\n Change your password through this url. you can just click the link to go directly:\n\n%s\n\nThank you,\nAffpilot AI Team",
		verificationURL,
	)

	err := email.SendVerificationEmail(req.Email, "Reset your Affpilot Password", mailText)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to send mail")
		return
	}

	utils.SuccessResponse(w, http.StatusAccepted, "password reset email has been sent", nil)
}

func ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req models.Password
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "token is required")
		return
	}

	cnf := config.GetConfig()
	secret := cnf.JWT.Secret

	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil || !parsedToken.Valid {
		log.Printf("Invalid token: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid token")
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid token claims")
		return
	}

	email, _ := claims["email"].(string)
	expiryFloat, _ := claims["expiry"].(float64)
	expiry := int64(expiryFloat)

	if time.Now().Unix() > expiry {
		utils.ErrorResponse(w, http.StatusRequestTimeout, "token has been expired")
		return
	}

	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)

	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Email is not registered")
		return
	}
	salt := os.Getenv("PASSWORD_SALT")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	_, err = database.DB.Exec(`
		UPDATE users
		SET password_hash = $1
		WHERE email =$2
	`, hashedPassword, email)

	if err != nil {
		log.Printf("Error updating password: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error (db)")
		return
	}

	utils.SuccessResponse(w, http.StatusAccepted, "password reset successful", nil)
}
