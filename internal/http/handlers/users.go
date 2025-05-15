package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
)

func Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userRole, err := middleware.GetUserRole(r)
	if err != err {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"Status":  strconv.Itoa(http.StatusUnauthorized),
			"message": "Log in first",
		})
		return
	}

	if userRole != "admin" && userRole != "system_admin" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"Status":  strconv.Itoa(http.StatusUnauthorized),
			"message": "Log in as admin",
		})
		return
	}

	rows, err := database.DB.Query(`
		SELECT username, email, first_name, last_name, email_verified, active
		FROM users
	`)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  strconv.Itoa(http.StatusAccepted),
		"message": "all users list",
		"data":    userlist,
	})
}
