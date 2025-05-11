package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
)

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func ListRoles(w http.ResponseWriter, r *http.Request) {
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
			"message": "Role checking only for admin+ user",
		})
		return
	}

	// Get all roles from database
	rows, err := database.DB.Query(`
        SELECT id, name, description 
        FROM roles 
        ORDER BY name ASC
    `)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			log.Printf("Error scanning role: %v", err)
			continue
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating roles: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"roles": roles,
		"count": len(roles),
	})
}
