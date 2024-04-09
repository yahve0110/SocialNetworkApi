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

// GetAllGroupRequestsHandler handles the retrieval of all group membership requests

func GetAllGroupRequestsHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the database package
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

	// Decode JSON payload from request body
	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the group ID from the decoded JSON payload
	groupID, ok := requestBody["group_id"]
	if !ok {
		http.Error(w, "Missing group_id in the request body", http.StatusBadRequest)
		return
	}

	// Check if the user is the creator of the group
	isGroupCreator, err := helpers.IsUserGroupCreator(dbConnection, userID, groupID)
	if err != nil {
		log.Printf("Error checking group creator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isGroupCreator {
		http.Error(w, "Unauthorized: Only group creator can view group requests", http.StatusUnauthorized)
		return
	}

	// Retrieve all group requests for the group
	groupRequests, err := getAllGroupRequestsFromDatabase(dbConnection, groupID)
	if err != nil {
		log.Printf("Error fetching group requests from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupRequests)
}

// getAllGroupRequestsFromDatabase fetches all group membership requests for a group from the database
func getAllGroupRequestsFromDatabase(dbConnection *sql.DB, groupID string) ([]models.GroupRequest, error) {
	// Query all group requests for the group from the "group_requests" table
	rows, err := dbConnection.Query("SELECT request_id, group_id, user_id, status FROM group_requests WHERE group_id = ?", groupID)
	if err != nil {
		log.Printf("Error querying group requests from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and create group request objects
	var groupRequests []models.GroupRequest
	for rows.Next() {
		var groupRequest models.GroupRequest
		if err := rows.Scan(&groupRequest.RequestID, &groupRequest.GroupID, &groupRequest.UserID, &groupRequest.Status); err != nil {
			log.Printf("Error scanning group request rows: %v", err)
			return nil, err
		}
		groupRequests = append(groupRequests, groupRequest)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over group request rows: %v", err)
		return nil, err
	}

	return groupRequests, nil
}
