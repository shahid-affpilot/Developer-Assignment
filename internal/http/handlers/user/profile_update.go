package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/database"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/models"
)

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIdFromParam, err := uuid.Parse(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	var update_info models.User
	if err := json.NewDecoder(r.Body).Decode(&update_info); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userIdFromParam).Scan(&exists)

	if err != nil {
		log.Printf("Database error checking user existence: %v", err)
		http.Error(w, "Internal server error1", http.StatusInternalServerError)
		return
	}

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusNotFound),
			"message": "User not found",
		})
		return
	}

	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", update_info.Username).Scan(&exists)
	// TODO username keep same as before

	if err != nil {
		log.Printf("Database error checking username existence: %v", err)
		http.Error(w, "Internal server error2", http.StatusInternalServerError)
		return
	}

	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusForbidden),
			"message": "Username already exists",
		})
		return
	}

	var user models.User
	err = database.DB.QueryRow(`
		UPDATE users
		SET username = $1,
		first_name = $2,
		last_name = $3
		WHERE id = $4
		RETURNING id, username, first_name, last_name
	`,
		update_info.Username,
		update_info.FirstName,
		update_info.LastName,
		userIdFromParam,
	).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
	)

	if err != nil {
		log.Printf("Error updating user detail: %v", err)
		http.Error(w, "Internal server error3", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  strconv.Itoa(http.StatusAccepted),
		"message": "user detail updated",
		"data":    user,
	})
}
