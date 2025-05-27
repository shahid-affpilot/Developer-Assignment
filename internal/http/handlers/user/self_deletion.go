package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
)

func UserDeletion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusUnauthorized),
			"message": "param is invalid",
		})
		return
	}

	fmt.Println("pera!!!")

	if userID != userIdFromParam.String() {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusUnauthorized),
			"message": "You can not delete other user",
		})
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusInternalServerError),
			"message": "Failed to update user",
		})
		return
	}

	fmt.Println("pera 3!!!")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "User deletion requested successfully",
	})
}
