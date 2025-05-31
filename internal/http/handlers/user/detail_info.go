package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func UserDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user id")
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
		utils.ErrorResponse(w, http.StatusNotFound, "user not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "user details", user)
}
