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

func UserDemote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userRole, err := middleware.GetUserRole(r)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusUnauthorized),
			"message": "log in first",
		})
		return
	}

	if userRole != "admin" && userRole != "system_admin" {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusUnauthorized),
			"message": "log in as admin+",
		})
		return
	}

	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "User ID in param is invalid",
		})
		return
	}

	var role RoleRequest
	json.NewDecoder(r.Body).Decode(&role)

	if role.RoleName == "admin" || role.RoleName == "system_admin" {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "this is not for promote user!!",
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

	var roleExists bool

	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)", role.RoleName).Scan(&roleExists)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusInternalServerError),
			"message": "Error checking role existence",
		})
		return
	}

	if !roleExists {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusNotFound),
			"message": "Role does not exist",
		})
		return
	}

	var role_id uuid.UUID
	err = database.DB.QueryRow("SELECT id FROM roles WHERE name = $1", role.RoleName).Scan(&role_id)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "database query problem2",
		})
		return
	}

	if !exists {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "there are no role in roles table",
		})
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
		"message": "user role demoted!",
	})
}
