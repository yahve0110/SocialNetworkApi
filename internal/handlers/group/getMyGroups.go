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

func GetMyGroups(w http.ResponseWriter, r *http.Request) {
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
	groups, err := getMyGroupsFromDatabase(dbConnection, userID)
	if err != nil {
		log.Printf("Error fetching groups from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// getAllGroupsFromDatabase fetches all groups from the database
func getMyGroupsFromDatabase(dbConnection *sql.DB, userID string) ([]models.Group, error) {
	// Query groups where the user is the creator
	// rows, err := dbConnection.Query("SELECT g.group_id, g.group_name, g.group_description, g.group_image, g.creation_date FROM groups g WHERE g.creator_id = $1", userID)
	// if err != nil {
	// 	log.Printf("Error querying groups created by the user: %v", err)
	// 	return nil, err
	// }
	// defer rows.Close()

	// Iterate through the result set and create group objects
	var groups []models.Group
	// for rows.Next() {
	// 	var group models.Group
	// 	if err := rows.Scan(&group.GroupID, &group.GroupName, &group.GroupDescription, &group.GroupImage, &group.CreationDate); err != nil {
	// 		log.Printf("Error scanning group rows: %v", err)
	// 		return nil, err
	// 	}
	// 	groups = append(groups, group)
	// }

	// // Check for errors from iterating over rows
	// if err := rows.Err(); err != nil {
	// 	log.Printf("Error iterating over group rows: %v", err)
	// 	return nil, err
	// }

	// Query groups where the user is a member
	rows, err := dbConnection.Query("SELECT g.group_id, g.group_name, g.group_description, g.group_image, g.creation_date FROM groups g JOIN group_members gm ON g.group_id = gm.group_id WHERE gm.user_id = $1", userID)
	if err != nil {
		log.Printf("Error querying groups where the user is a member: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and add group objects
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.GroupID, &group.GroupName, &group.GroupDescription, &group.GroupImage, &group.CreationDate); err != nil {
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
