package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
)

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"required"`
}

func CreateRole(w http.ResponseWriter, r *http.Request) {
	// Check user role
	userRole, err := middleware.GetUserRole(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verify if user is admin or system_admin
	if userRole != "admin" && userRole != "system_admin" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Role creation only for admin+ user",
		})
		return
	}

	// Parse request body
	var req CreateRoleRequest
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

	// Check if role name already exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)", req.Name).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking role existence: %v", err)
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

	// Create role in database
	var role RoleDetails
	err = database.DB.QueryRow(`
        INSERT INTO roles (name, description)
        VALUES ($1, $2)
        RETURNING id, name, description, created_at, updated_at`,
		req.Name,
		req.Description,
	).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		log.Printf("Error creating role: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}
