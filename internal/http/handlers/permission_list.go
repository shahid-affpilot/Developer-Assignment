package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
)

func GetPermissionList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := middleware.GetUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Please login first",
		})
		return
	}
	if userID != "admin" && userID != "system_admin" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Only admin users can view permissions",
		})
		return
	}

	var permissions []Permission
	rows, err := database.DB.Query(`
        SELECT DISTINCT
            p.id,
            p.name,
            p.description,
            p.resource,
            p.action,
            p.created_at,
            p.updated_at
        FROM permissions p
        ORDER BY p.name`)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to fetch permissions",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var perm Permission
		err := rows.Scan(
			&perm.ID,
			&perm.Name,
			&perm.Description,
			&perm.Resource,
			&perm.Action,
			&perm.CreatedAt,
			&perm.UpdatedAt,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Error scanning permissions",
			})
			return
		}
		permissions = append(permissions, perm)
	}

	if err = rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error iterating permissions",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"count":  len(permissions),
		"data":   permissions,
	})
}
