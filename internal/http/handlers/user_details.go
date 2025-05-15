package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
)

func UserDetail(w http.ResponseWriter, r *http.Request) {
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

	userID, _ := middleware.GetUserID(r)
	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	if userRole != "moderator" && userRole != "admin" && userRole != "system_admin" && userID != userIdFromParam.String() {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"Status":  strconv.Itoa(http.StatusUnauthorized),
			"message": "Log in as admin",
		})
		return
	}

	var user models.User
	err = database.DB.QueryRow(`
		SELECT username, email, first_name, last_name, email_verified, active, user_type
		FROM users
		WHERE id = $1`,
		userIdFromParam.String(),
	).Scan(&user.Username, &user.Email, &user.FirstName, &user.LastName, &user.EmailVerified, &user.Active, &user.UserType)

	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusNotFound),
			"message": "User not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": strconv.Itoa(http.StatusAccepted),
		"data":   user,
	})
}
