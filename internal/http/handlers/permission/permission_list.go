package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
)

func GetPermissionList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var permissions []models.PermissionShort
	rows, err := database.DB.Query(`
        SELECT DISTINCT
            p.id,
            p.name,
            p.description
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
		var perm models.PermissionShort
		err := rows.Scan(
			&perm.ID,
			&perm.Name,
			&perm.Description,
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
		"status":  strconv.Itoa(http.StatusOK),
		"message": "permission details",
		"data":    permissions,
	})
}
