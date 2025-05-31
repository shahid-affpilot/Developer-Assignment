package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var update_info models.User
	if err := json.NewDecoder(r.Body).Decode(&update_info); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userIdFromParam).Scan(&exists)

	if err != nil {
		log.Printf("Database error checking user existence: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if !exists {
		utils.ErrorResponse(w, http.StatusNotFound, "user not found")
		return
	}

	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND id != $2)", update_info.Username, userIdFromParam).Scan(&exists)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if exists {
		utils.ErrorResponse(w, http.StatusConflict, "username already exists")
		return
	}

	var user models.User
	err = database.DB.QueryRow(`
		UPDATE users
		SET username = $1,
		first_name = $2,
		last_name = $3
		WHERE id = $4
		RETURNING id, username, first_name, last_name
	`,
		update_info.Username,
		update_info.FirstName,
		update_info.LastName,
		userIdFromParam,
	).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
	)

	if err != nil {
		log.Printf("Error updating user detail: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "User updated", user)
}
