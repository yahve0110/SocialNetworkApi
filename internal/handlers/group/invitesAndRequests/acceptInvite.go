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

// AcceptGroupInvitationHandler handles the acceptance of group invitations

func AcceptGroupInvitationHandler(w http.ResponseWriter, r *http.Request) {
	var invitationData models.GroupInvitation

	if err := json.NewDecoder(r.Body).Decode(&invitationData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Check if the requester is the user receiving the invitation
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

	// Check if the invitation exists
	exists, err := helpers.InvitationExists(dbConnection, invitationData.GroupID, userID)
	if err != nil {
		log.Printf("Error checking if invitation exists: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Invitation not found or not intended for the current user", http.StatusBadRequest)
		return
	}

	// Add the user to the group members
	err = AddUserToGroup(dbConnection, userID, invitationData.GroupID)
	if err != nil {
		log.Printf("Error adding user to group: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = helpers.DeleteInvitation(dbConnection, invitationData.GroupID, userID)
	if err != nil {
		log.Printf("Error deleting invitation: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invitationData)
}

// AddUserToGroup adds a user to a group
func AddUserToGroup(dbConnection *sql.DB, userID, groupID string) error {
	// This is a simplified example assuming a SQL database with a "group_members" table

	// Check if the user is already a member of the group
	isMember, err := helpers.IsUserGroupMember(dbConnection, userID, groupID)
	if err != nil {
		log.Printf("Error checking if user is already a group member: %v", err)
		return err
	}
	if isMember {
		log.Printf("User is already a member of the group")
		return nil
	}

	// Add the user to the group
	_, err = dbConnection.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, userID)
	if err != nil {
		log.Printf("Error adding user to group in database: %v", err)
		return err
	}

	return nil
}
