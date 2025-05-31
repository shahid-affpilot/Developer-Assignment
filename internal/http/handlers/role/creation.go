package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func CreateRole(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req models.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate request
	if req.Name == "" || len(req.Name) < 2 || len(req.Name) > 50 {
		utils.ErrorResponse(w, http.StatusBadRequest, "role name must be in between 3-50 char")
		return
	}

	// Check if role name already exists
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)", req.Name).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking role existence: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error (db)")
		return
	}

	if exists {
		w.Header().Set("Content-Type", "application/json")
		utils.ErrorResponse(w, http.StatusConflict, "role (name) already exists")
		return
	}

	// Create role in database
	var role models.RoleDetails
	err = database.DB.QueryRow(`
        INSERT INTO roles (name, description)
        VALUES ($1, $2)
        RETURNING id, name, description, created_at, updated_at`,
		req.Name,
		req.Description,
	).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		log.Printf("Error creating role: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "role creation success", role)
}
