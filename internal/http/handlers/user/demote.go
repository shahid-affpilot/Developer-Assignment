package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func UserDemote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var role models.RoleRequest
	json.NewDecoder(r.Body).Decode(&role)

	if role.RoleName == "admin" || role.RoleName == "system_admin" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "can not set this role")
		return
	}

	var exists bool
	err = database.DB.QueryRow("SELECT exists(SELECT 1 FROM users WHERE id=$1)", userIdFromParam).Scan(&exists)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !exists {
		utils.ErrorResponse(w, http.StatusNoContent, "no user found with this id")
		return
	}

	var roleExists bool

	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)", role.RoleName).Scan(&roleExists)
	if err != nil {
		utils.ErrorResponse(w, http.StatusForbidden, "error checking role")
		return
	}

	if !roleExists {
		utils.ErrorResponse(w, http.StatusNoContent, "role did not found")
		return
	}

	var role_id uuid.UUID
	err = database.DB.QueryRow("SELECT id FROM roles WHERE name = $1", role.RoleName).Scan(&role_id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !exists {
		utils.ErrorResponse(w, http.StatusNotFound, "role did not found")
		return
	}

	_, err = database.DB.Exec(`
		UPDATE users
		SET user_type = $1
		WHERE id = $2
	`,
		role.RoleName,
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

	utils.SuccessResponse(w, http.StatusOK, "user role demoted", nil)
}
