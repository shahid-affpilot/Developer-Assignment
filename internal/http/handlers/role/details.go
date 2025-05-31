package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func GetRole(w http.ResponseWriter, r *http.Request) {
	// Get role ID from URL params
	vars := mux.Vars(r)
	roleID, err := uuid.Parse(vars["role_id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid role id")
		return
	}

	// Get role details from database
	var role models.RoleDetails
	err = database.DB.QueryRow(`
        SELECT id, name, description, created_at, updated_at 
        FROM roles 
        WHERE id = $1`,
		roleID,
	).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)

	if err == sql.ErrNoRows {
		utils.ErrorResponse(w, http.StatusBadRequest, "role not found")
		return
	}

	if err != nil {
		log.Printf("Database error: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "user information detail", role)
}
