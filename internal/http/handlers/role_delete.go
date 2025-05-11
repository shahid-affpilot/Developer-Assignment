package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
)

func DeleteRole(w http.ResponseWriter, r *http.Request) {
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
			"message": "Role deletion only for admin+ user",
		})
		return
	}

	// Get role ID from URL params
	vars := mux.Vars(r)
	roleID, err := uuid.Parse(vars["role_id"])
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if isSystemRole {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Cannot delete system roles",
		})
		return
	}

	// Delete role
	result, err := database.DB.Exec("DELETE FROM roles WHERE id = $1", roleID)
	if err != nil {
		log.Printf("Error deleting role: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// how many row affected or deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Role not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Role deleted successfully",
	})
}
