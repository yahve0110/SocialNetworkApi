package groupHandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	var groupData models.Group

	if err := json.NewDecoder(r.Body).Decode(&groupData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate that required fields are not empty
	if groupData.GroupName == "" || groupData.GroupDescription == "" {
		http.Error(w, "Group Name, and GroupDescription are required fields", http.StatusBadRequest)
		return
	}

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

	// Set the author ID for the new group
	groupData.CreatorID = userID
	// Generate a UUID for GroupID
	groupData.GroupID = uuid.New().String()
	groupData.CreationDate = time.Now()

	// Attempt to save the group to the database
	err = saveGroupToDatabase(groupData, dbConnection)
	if err != nil {
		// Check if the error is due to non-unique group name
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, fmt.Sprintf("Group with name '%s' already exists", groupData.GroupName), http.StatusBadRequest)
			return
		}

		log.Printf("Error saving group to database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groupData); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// InsertGroup inserts a new group into the database
func saveGroupToDatabase(groupData models.Group, db *sql.DB) error {
	// Check if the group name is unique
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM groups WHERE group_name = ?", groupData.GroupName).Scan(&count)
	if err != nil {
		log.Printf("Error checking group name uniqueness: %v", err)
		return err
	}

	// If count is greater than 0, a group with the same name already exists
	if count > 0 {
		return fmt.Errorf("group with name '%s' already exists", groupData.GroupName)
	}

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO groups (group_id, group_name, group_description, creator_id, creation_date, group_image) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(groupData.GroupID, groupData.GroupName, groupData.GroupDescription, groupData.CreatorID, groupData.CreationDate, groupData.GroupImage)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}
