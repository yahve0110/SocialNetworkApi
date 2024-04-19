package groupHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/models"
)

func GetGroupByID(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Get the group ID from the URL query parameters
	groupID := r.URL.Query().Get("groupID")
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Query the database to fetch the group by ID
	group, err := getGroupByIDFromDatabase(dbConnection, groupID)
	if err != nil {
		log.Printf("Error fetching group by ID from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if the group exists
	if group.GroupID == "" {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	// Serialize the group object to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}

// getGroupByIDFromDatabase fetches the group from the database by ID
func getGroupByIDFromDatabase(dbConnection *sql.DB, groupID string) (models.Group, error) {
	var group models.Group

	// Query the database to fetch the group by ID
	err := dbConnection.QueryRow("SELECT group_id, group_name, group_description, group_image, creation_date, creator_id FROM groups WHERE group_id = $1", groupID).
		Scan(&group.GroupID, &group.GroupName, &group.GroupDescription, &group.GroupImage, &group.CreationDate, &group.CreatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows were returned, it means the group does not exist
			return group, nil
		}
		// Other errors
		log.Printf("Error querying group by ID: %v", err)
		return group, err
	}

	return group, nil
}
