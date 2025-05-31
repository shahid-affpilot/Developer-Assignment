package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/config"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "token is missing in URL")
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
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid token")
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {

		utils.ErrorResponse(w, http.StatusUnauthorized, "invalid token claims")
		return
	}

	userIDStr, _ := claims["user_id"].(string)
	expiryFloat, _ := claims["expiry"].(float64)
	expiry := int64(expiryFloat)

	if time.Now().Unix() > expiry {
		utils.ErrorResponse(w, http.StatusBadRequest, "token has been expired")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Failed to parsed user id")
		return
	}

	_, err = database.DB.Exec(`
		UPDATE users 
		SET email_verified = true,
		    verification_token = NULL,
		    token_expiry = NULL
		WHERE id = $1
	`, userID)

	if err != nil {
		log.Printf("Error updating verification status: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Internal server error, failed to update user")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "user is now verified", nil)
}
