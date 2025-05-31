package handlers

import (
	"net/http"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
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
		utils.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
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
			utils.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		permissions = append(permissions, perm)
	}

	if err = rows.Err(); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Permission list", permissions)
}
