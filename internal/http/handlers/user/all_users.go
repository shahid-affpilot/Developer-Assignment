package handlers

import (
	"log"
	"net/http"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := database.DB.Query(`
		SELECT username, email, first_name, last_name, email_verified, active
		FROM users
	`)
	if err != nil {
		log.Printf("Database error: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer rows.Close()

	var userlist []models.UserShortInfo
	for rows.Next() {
		var user models.UserShortInfo
		if err := rows.Scan(&user.UserName, &user.Email, &user.FirstName, &user.LastName, &user.EmailVerified, &user.Active); err != nil {
			log.Printf("Error scanning role: %v", err)
			continue
		}
		userlist = append(userlist, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating roles: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "all user list", userlist)
}
