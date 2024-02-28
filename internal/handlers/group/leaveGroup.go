package groupHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

type LeaveGroupRequest struct {
	GroupID string `json:"group_id"`
}

// LeaveGroupHandler handles the user leaving a group

func LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	var requestData LeaveGroupRequest

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

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

	// Check if the user is a member of the group
	isGroupMember, err := helpers.IsUserGroupMember(dbConnection, userID, requestData.GroupID)
	if err != nil {
		log.Printf("Error checking group membership: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isGroupMember {
		http.Error(w, "Unauthorized: User is not a member of the group", http.StatusUnauthorized)
		return
	}

	// Remove the user from the group
	err = RemoveUserFromGroup(dbConnection, userID, requestData.GroupID)
	if err != nil {
		log.Printf("Error removing user from group: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully left the group"})
}

// RemoveUserFromGroup removes a user from a group
func RemoveUserFromGroup(db *sql.DB, userID, groupID string) error {
	// Remove the user from the group
	_, err := db.Exec("DELETE FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID)
	if err != nil {
		log.Printf("Error removing user from group in database: %v", err)
		return err
	}

	return nil
}
