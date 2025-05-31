package handlers

import (
	"net/http"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	var permissions []models.PermissionShort
	rows, err := database.DB.Query(`
        SELECT DISTINCT
            p.id,
            p.name,
            p.description
        FROM permissions p
        INNER JOIN role_permissions rp ON p.id = rp.permission_id
        INNER JOIN user_roles ur ON rp.role_id = ur.role_id
        WHERE ur.user_id = $1
        ORDER BY p.name`, userID)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
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
			utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
			return
		}
		permissions = append(permissions, perm)
	}

	utils.SuccessResponse(w, http.StatusOK, "success", permissions)
}
