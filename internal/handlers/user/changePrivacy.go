package userHandlers

import (
	"database/sql"
	"encoding/json"

	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

type UpdatePrivacyRequest struct {
	Privacy string `json:"privacy"`
}

func UpdatePrivacy(w http.ResponseWriter, r *http.Request) {
	// Get the database connection
	dbConnection := database.DB

	// Parse request body
	var updateRequest UpdatePrivacyRequest
	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the user ID based on the current user's session
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get the user ID based on the current user's session
	userID, err := helpers.GetUserIDFromSession(dbConnection, cookie.Value)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Update user privacy in the database
	err = UpdateUserPrivacy(dbConnection, userID, updateRequest.Privacy)
	if err != nil {
		http.Error(w, "Failed to update user privacy", http.StatusInternalServerError)
		return
	}

	// Respond with success message and updated privacy
	response := map[string]string{
		"message": "Privacy updated successfully",
		"privacy": updateRequest.Privacy,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUserPrivacy(db *sql.DB, userID, privacy string) error {
	stmt, err := db.Prepare("UPDATE users SET privacy = ? WHERE user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(privacy, userID)
	if err != nil {
		return err
	}

	return nil
}
