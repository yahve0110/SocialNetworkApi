package groupInviteHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

// GetUserInvitationsHandler handles the retrieval of group invitations for a specific user
func GetUserInvitationsHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Get the user ID based on the current user's session
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

	// Retrieve group invitations for the user
	invitations, err := helpers.GetUserInvitations(dbConnection, userID)
	if err != nil {
		log.Printf("Error fetching user invitations from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Enhance the invitations with group names
	enhancedInvitations, err := enhanceInvitationsWithGroupNames(dbConnection, invitations)
	if err != nil {
		log.Printf("Error enhancing invitations with group names: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enhancedInvitations)
}

// enhanceInvitationsWithGroupNames retrieves group names and enhances the invitations with them
func enhanceInvitationsWithGroupNames(dbConnection *sql.DB, invitations []models.GroupInvitation) ([]models.GroupInvitation, error) {
	// Prepare a map to store group IDs and their corresponding names
	groupNames := make(map[string]string)

	// Retrieve group names for each unique group ID
	for _, invitation := range invitations {
		// Check if the group ID is already in the map
		if _, ok := groupNames[invitation.GroupID]; !ok {
			// Retrieve and store the group name for the group ID
			groupName, err := GetGroupNameByID(dbConnection, invitation.GroupID)
			if err != nil {
				return nil, err
			}
			groupNames[invitation.GroupID] = groupName
		}
	}

	// Enhance the invitations with group names
	for i, invitation := range invitations {
		invitations[i].GroupName = groupNames[invitation.GroupID]
	}

	return invitations, nil
}

// GetGroupNameByID retrieves the name of a group by its ID from the database
func GetGroupNameByID(dbConnection *sql.DB, groupID string) (string, error) {
	// Query to retrieve the group name by its ID
	query := "SELECT group_name FROM groups WHERE group_id = ?"

	// Execute the query
	row := dbConnection.QueryRow(query, groupID)

	// Variable to store the group name
	var groupName string

	// Scan the result into the groupName variable
	err := row.Scan(&groupName)
	if err != nil {
		log.Printf("Error scanning group name: %v", err)
		return "", err
	}

	return groupName, nil
}
