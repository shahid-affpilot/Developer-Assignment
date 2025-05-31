package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func UserDeletion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])

	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	fmt.Println("pera!!!")

	if userID != userIdFromParam.String() {
		utils.ErrorResponse(w, http.StatusUnauthorized, "you cannot delete other user")
		return
	}

	fmt.Println("pera 2!!!")

	err = database.DB.QueryRow(`
		UPDATE users
		SET deletion_requested = TRUE,
		active = FALSE
		WHERE id = $1
		RETURNING id
	`,
		userID,
	).Scan(&userID)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	fmt.Println("pera 3!!!")

	utils.SuccessResponse(w, http.StatusOK, "deletion request submitted", nil)
}
