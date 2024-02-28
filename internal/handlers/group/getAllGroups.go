package groupHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/models"
)

func GetAllGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// This is a simplified example assuming a SQL database with a "groups" table
	groups, err := getAllGroupsFromDatabase(dbConnection)
	if err != nil {
		log.Printf("Error fetching groups from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// getAllGroupsFromDatabase fetches all groups from the database
func getAllGroupsFromDatabase(dbConnection *sql.DB) ([]models.Group, error) {
	// Query all groups from the "groups" table
	rows, err := dbConnection.Query("SELECT group_id, group_name, group_description, creator_id, creation_date FROM groups")
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
