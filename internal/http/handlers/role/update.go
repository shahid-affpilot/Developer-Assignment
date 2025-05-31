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

type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"required"`
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Get role ID from URL params
	vars := mux.Vars(r)
	roleID, err := uuid.Parse(vars["role_id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid role id")
		return
	}

	// Parse request body
	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate request
	if req.Name == "" || len(req.Name) < 2 || len(req.Name) > 50 {
		utils.ErrorResponse(w, http.StatusBadRequest, "Name must be between 2 and 50 characters")
		return
	}

	// Check if role exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM roles WHERE id = $1)", roleID).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking role existence: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !exists {
		utils.ErrorResponse(w, http.StatusNoContent, "no role with this id")
		return
	}

	// Check if new name already exists for different role
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1 AND id != $2)",
		req.Name, roleID,
	).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking name existence: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error (db)")
		return
	}

	if exists {
		utils.ErrorResponse(w, http.StatusConflict, "role with this name already exists")
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
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Role updated successfully", updatedRole)
}
