package groupHandlers

import (
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

// CreateGroupInvitationHandler handles the creation of group invitations

func SendGroupInvitationHandler(w http.ResponseWriter, r *http.Request) {
	var invitationData models.GroupInvitation

	if err := json.NewDecoder(r.Body).Decode(&invitationData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Check if the requester is the group creator or a member
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

	// Check if the user is the group creator or a member
	isGroupMember, err := helpers.IsUserGroupMember(dbConnection, userID, invitationData.GroupID)
	if err != nil {
		log.Printf("Error checking group membership: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	isGroupCreator, err := helpers.IsUserGroupCreator(dbConnection, userID, invitationData.GroupID)
	if err != nil {
		log.Printf("Error checking group creator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isGroupMember && !isGroupCreator {
		http.Error(w, "Unauthorized: Only group creator or members can send invitations", http.StatusUnauthorized)
		return
	}

	// Check if the receiver user exists
	exists, err := helpers.UserExists(dbConnection, invitationData.ReceiverID)
	if err != nil {
		http.Error(w, "Error checking if receiver user exists", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Receiver user not found", http.StatusBadRequest)
		return
	}

	// Check if an invitation from the same sender to the same receiver for the same group already exists
	if exists, err := helpers.InvitationExists(dbConnection, invitationData.GroupID, userID); err != nil {
		http.Error(w, "Error checking existing invitation", http.StatusInternalServerError)
		return
	} else if exists {
		http.Error(w, "Invitation already sent to the receiver", http.StatusBadRequest)
		return
	}

	// Set the inviter ID for the invitation
	invitationData.InviterID = userID

	//Set invitation status
	invitationData.Status = "pending"

	err = helpers.SaveInvitationToDatabase(invitationData, dbConnection)
	if err != nil {
		log.Printf("Error saving invitation to database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(invitationData); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
