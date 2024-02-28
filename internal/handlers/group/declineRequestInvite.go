package groupHandlers

import (
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

// RefuseGroupInvitationHandler handles the refusal of group invitations

func RefuseGroupInvitationHandler(w http.ResponseWriter, r *http.Request) {
	var invitationData models.GroupInvitation

	if err := json.NewDecoder(r.Body).Decode(&invitationData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the database package
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
	exists, err := helpers.InvitationExists(dbConnection, invitationData.GroupID,  userID)
	if err != nil {
		log.Printf("Error checking if invitation exists: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Invitation not found", http.StatusBadRequest)
		return
	}



	err = helpers.DeleteInvitation(dbConnection,invitationData.GroupID,userID)
	if err != nil {
		log.Printf("Error deleting invitation: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invitationData)
}
