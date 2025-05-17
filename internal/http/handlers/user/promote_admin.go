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

func PromoteToAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "User ID in param is invalid",
		})
		return
	}

	var exists bool
	err = database.DB.QueryRow("SELECT exists(SELECT 1 FROM users WHERE id=$1)", userIdFromParam).Scan(&exists)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "database query problem1",
		})
		return
	}

	if !exists {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "there are no user with this ID",
		})
		return
	}

	var role_id uuid.UUID
	err = database.DB.QueryRow("SELECT id FROM roles WHERE name = $1", "admin").Scan(&role_id)
	if err != nil {
		fmt.Printf("chillErr: %s", err)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "database query problem2",
		})
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
		fmt.Printf("chillErr: %s", err)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "database query problem3",
		})
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
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "database query problem4",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "promoted to admin!",
	})
}
