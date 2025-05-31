package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func PromoteToAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var exists bool
	err = database.DB.QueryRow("SELECT exists(SELECT 1 FROM users WHERE id=$1)", userIdFromParam).Scan(&exists)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !exists {
		utils.ErrorResponse(w, http.StatusNoContent, "there are no user with this id")
		return
	}

	var role_id uuid.UUID
	err = database.DB.QueryRow("SELECT id FROM roles WHERE name = $1", "admin").Scan(&role_id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error (db)")
		return
	}

	_, err = database.DB.Exec(`
		UPDATE users
		SET user_type = $1
		WHERE id = $2
	`,
		"admin",
		userIdFromParam,
	)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	userID, _ := middleware.GetUserID(r) // assigned by this user

	_, err = database.DB.Exec(`
		UPDATE user_roles
		SET role_id = $1,
		assigned_by = $2
		WHERE user_id = $3
	`,
		role_id,
		userID,
		userIdFromParam,
	)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "promoted to admin", nil)
}
