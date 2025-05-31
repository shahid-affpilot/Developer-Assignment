package handlers

import (
	"log"
	"net/http"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func ListRoles(w http.ResponseWriter, r *http.Request) {
	// Get all roles from database
	rows, err := database.DB.Query(`
        SELECT id, name, description 
        FROM roles 
        ORDER BY name ASC
    `)
	if err != nil {
		log.Printf("Database error: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			log.Printf("Error scanning role: %v", err)
			continue
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating roles: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SuccessResponse(w, http.StatusAccepted, "list of user roles", roles)
}
