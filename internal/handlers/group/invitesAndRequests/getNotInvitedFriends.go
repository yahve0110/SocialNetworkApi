package groupInviteHandlers

import (
	"database/sql"
	"encoding/json"

	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

// GetUninvitedFollowersHandler handles the retrieval of uninvited followers for a group
func GetUninvitedFollowersHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Check if the user is authenticated
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	userID, err := helpers.GetUserIDFromSession(dbConnection, cookie.Value)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get the group ID from the query parameters
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		http.Error(w, "Missing group_id in the request URL", http.StatusBadRequest)
		return
	}

	// Retrieve uninvited followers for the group
	uninvitedFollowers, err := getUninvitedFollowers(dbConnection, userID, groupID)
	if err != nil {
		log.Printf("Error fetching uninvited followers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send the uninvited followers as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uninvitedFollowers)
}

// getUninvitedFollowers fetches uninvited followers for a group from the database
func getUninvitedFollowers(dbConnection *sql.DB, userID, groupID string) ([]map[string]interface{}, error) {
	// Query uninvited followers for the group from the database
	query := `
	SELECT u.user_id, u.first_name, u.last_name, u.profile_picture
	FROM users u
	WHERE u.user_id != ?
	AND u.user_id NOT IN (
		SELECT receiver_id
		FROM group_invitations
		WHERE group_id = ?
	)
	AND u.user_id NOT IN (
		SELECT user_id
		FROM group_members
		WHERE group_id = ?
	)
	AND u.user_id != (
		SELECT creator_id
		FROM groups
		WHERE group_id = ?
	)
    `

	rows, err := dbConnection.Query(query, userID, groupID, groupID, groupID)
	if err != nil {
		log.Printf("Error executing SQL query: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Create a slice to store the uninvited followers
	var uninvitedFollowers []map[string]interface{}

	// Iterate over the query results and add them to the slice
	for rows.Next() {
		var userID, firstName, lastName, profilePicture string
		if err := rows.Scan(&userID, &firstName, &lastName, &profilePicture); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		uninvitedFollower := map[string]interface{}{
			"user_id":        userID,
			"first_name":     firstName,
			"last_name":      lastName,
			"profilePicture": profilePicture,
		}

		uninvitedFollowers = append(uninvitedFollowers, uninvitedFollower)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return uninvitedFollowers, nil
}

