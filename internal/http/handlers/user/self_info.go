package handlers

import (
	"fmt"
	"net/http"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/http/middleware"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func UserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	// Get user information from database
	var user models.UserResponse
	err := database.DB.QueryRow(`
        SELECT id, username, first_name, last_name, email, user_type, email_verified, created_at, updated_at
        FROM users
        WHERE id = $1`,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.Verified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	fmt.Println(err)

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Profile info", user)
}
