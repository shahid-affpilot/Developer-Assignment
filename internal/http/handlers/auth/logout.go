package handlers

import (
	"net/http"
	"time"

	"github.com/shahid-affpilot/affpilot-auth-service/internal/utils"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
	})

	utils.SuccessResponse(w, http.StatusOK, "user logged out success", nil)
}
