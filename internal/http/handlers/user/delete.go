package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func UserDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])

	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	err = database.DB.QueryRow(`
		UPDATE users
		SET deletion_requested = TRUE,
		active = FALSE
		WHERE id = $1
		RETURNING id
	`,
		userIdFromParam,
	).Scan(&userIdFromParam)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "user deleted", nil)
}
