package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
)

func PermissionDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userRole, err := middleware.GetUserRole(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Please login first",
		})
		return
	}
	if userRole != "admin" && userRole != "system_admin" {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "be a admit first",
		})
		return
	}

	var permission models.PermissionDetails

	userID, _ := middleware.GetUserID(r)

	err = database.DB.QueryRow(`
    SELECT
        p.id,
        p.name,
        p.description,
        p.resource,
        p.action,
        p.created_at,
        p.updated_at
    FROM permissions p
    INNER JOIN role_permissions rp ON p.id = rp.permission_id
    INNER JOIN user_roles ur ON rp.role_id = ur.role_id
    WHERE ur.user_id = $1
    ORDER BY p.name
    LIMIT 1`, userID).Scan(
		&permission.ID,
		&permission.Name,
		&permission.Description,
		&permission.Resource,
		&permission.Action,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No permission found for this user.")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "No permission found for this user",
			})
			return
		} else {
			log.Printf("Error querying permission: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Error querying permission",
				"error":   err.Error(),
			})
			return
		}
	}

	fmt.Printf("permission list: %+v", permission)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   permission,
	})
}
