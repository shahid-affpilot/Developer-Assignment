package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
)

type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"required"`
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Get role ID from URL params
	vars := mux.Vars(r)
	roleID, err := uuid.Parse(vars["role_id"])
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" || len(req.Name) < 2 || len(req.Name) > 50 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Name must be between 2 and 50 characters",
		})
		return
	}

	// Check if role exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM roles WHERE id = $1)", roleID).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking role existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Role not found",
		})
		return
	}

	// Check if new name already exists for different role
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1 AND id != $2)",
		req.Name, roleID,
	).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking name existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Role with this name already exists",
		})
		return
	}

	// Update role in database
	var updatedRole models.RoleDetails
	err = database.DB.QueryRow(`
        UPDATE roles 
        SET name = $1, 
            description = $2,
            updated_at = NOW()
        WHERE id = $3
        RETURNING id, name, description, created_at, updated_at`,
		req.Name,
		req.Description,
		roleID,
	).Scan(
		&updatedRole.ID,
		&updatedRole.Name,
		&updatedRole.Description,
		&updatedRole.CreatedAt,
		&updatedRole.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error updating role: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  strconv.Itoa(http.StatusAccepted),
		"message": "role update successful",
		"data":    updatedRole,
	})
}
