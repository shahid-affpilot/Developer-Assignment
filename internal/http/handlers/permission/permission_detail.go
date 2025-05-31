package handlers

import (
	"database/sql"
	"net/http"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func PermissionDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var permission models.PermissionDetails

	userID, _ := middleware.GetUserID(r)

	err := database.DB.QueryRow(`
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
			utils.ErrorResponse(w, http.StatusNoContent, "No permission for this user")
			return
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	utils.SuccessResponse(w, http.StatusOK, "List of permissions", permission)
}
