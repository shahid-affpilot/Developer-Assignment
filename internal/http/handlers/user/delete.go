package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
)

func UserDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "user id invalid",
		})
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
		fmt.Printf("err: %s", err)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "error to delete user",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "User deletion requested successfully",
	})
}
