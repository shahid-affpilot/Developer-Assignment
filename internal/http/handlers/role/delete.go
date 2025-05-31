package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Get role ID from URL params
	vars := mux.Vars(r)
	roleID, err := uuid.Parse(vars["role_id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid role id")
		return
	}

	// Check if role exists and is not a system role
	var isSystemRole bool
	err = database.DB.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM roles 
            WHERE id = $1 AND name IN ('system_admin', 'admin', 'moderator', 'user')
        )`, roleID).Scan(&isSystemRole)

	if err != nil {
		log.Printf("Database error checking role: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if isSystemRole {
		utils.ErrorResponse(w, http.StatusBadRequest, "system role are default, cant be deleted")
		return
	}

	// Delete role
	result, err := database.DB.Exec("DELETE FROM roles WHERE id = $1", roleID)
	if err != nil {
		log.Printf("Error deleting role: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// how many row affected or deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if rowsAffected == 0 {
		utils.ErrorResponse(w, http.StatusNoContent, "role not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Role deleted", nil)
}
