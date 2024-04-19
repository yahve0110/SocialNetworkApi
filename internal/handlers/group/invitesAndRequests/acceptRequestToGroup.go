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

// AcceptGroupRequestHandler handles the acceptance of group membership requests

func AcceptGroupRequestHandler(w http.ResponseWriter, r *http.Request) {
	var requestData models.GroupRequest

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Check if the requester is the group creator
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	creatorID, err := helpers.GetUserIDFromSession(dbConnection, cookie.Value)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Check if the user is the group creator
	isGroupCreator, err := helpers.IsUserGroupCreator(dbConnection, creatorID, requestData.GroupID)
	if err != nil {
		log.Printf("Error checking group creator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isGroupCreator {
		http.Error(w, "Unauthorized: Only group creator can accept membership requests", http.StatusUnauthorized)
		return
	}

	// Check if the request exists
	exists, err := helpers.GroupRequestExists(dbConnection, requestData.GroupID, requestData.UserID)
	if err != nil {
		log.Printf("Error checking if group request exists: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Group request not found", http.StatusBadRequest)
		return
	}

	// Add the user to the group members
	err = AddUserToGroup(dbConnection, requestData.UserID, requestData.GroupID)
	if err != nil {
		log.Printf("Error adding user to group: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Delete the group request
	err = DeleteGroupRequest(dbConnection, requestData.GroupID, requestData.UserID)
	if err != nil {
		log.Printf("Error deleting group request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set response headers and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requestData)
}

// DeleteGroupRequest deletes a group request from the database
func DeleteGroupRequest(db *sql.DB, GroupID, UserID string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("DELETE FROM group_requests WHERE group_id = ? AND user_id = ?")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(GroupID, UserID)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}
