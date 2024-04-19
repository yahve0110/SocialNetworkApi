package groupHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

func GetAllGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

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

	// This is a simplified example assuming a SQL database with a "groups" table
	groups, err := getAllGroupsFromDatabase(dbConnection, userID)
	if err != nil {
		log.Printf("Error fetching groups from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// getAllGroupsFromDatabase fetches all groups from the database excluding those the user is a member or owner of
func getAllGroupsFromDatabase(dbConnection *sql.DB, userID string) ([]models.Group, error) {
	// Query groups from the "groups" table excluding those where the user is a member, creator, or has sent a request
	rows, err := dbConnection.Query(`
		SELECT g.group_id, g.group_name, g.group_description, g.creator_id, g.creation_date
		FROM groups g
		WHERE g.creator_id != $1
		AND NOT EXISTS (
			SELECT 1 FROM group_members gm WHERE gm.group_id = g.group_id AND gm.user_id = $1
		)
		AND NOT EXISTS (
			SELECT 1 FROM group_requests gr WHERE gr.group_id = g.group_id AND gr.user_id = $1
		)
	`, userID)
	if err != nil {
		log.Printf("Error querying groups from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and create group objects
	var groups []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.GroupID, &group.GroupName, &group.GroupDescription, &group.CreatorID, &group.CreationDate); err != nil {
			log.Printf("Error scanning group rows: %v", err)
			return nil, err
		}
		groups = append(groups, group)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over group rows: %v", err)
		return nil, err
	}

	return groups, nil
}

