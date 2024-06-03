package groupHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

// GroupRequestWithName represents a group request with the group name
type GroupRequestWithName struct {
	RequestID string `json:"request_id"`
	GroupID   string `json:"group_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	GroupName string `json:"group_name"`
}

func GetRequestedGroups(w http.ResponseWriter, r *http.Request) {
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

	// This is a simplified example assuming a SQL database with a "group_requests" table
	groupRequests, err := getGroupRequestsSentFromDatabase(dbConnection, userID)
	if err != nil {
		log.Printf("Error fetching group requests sent by user from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupRequests)
}

// getGroupRequestsSentFromDatabase fetches all group requests sent by the user with group names
func getGroupRequestsSentFromDatabase(dbConnection *sql.DB, userID string) ([]GroupRequestWithName, error) {
	// Query group requests from the "group_requests" table sent by the user with group names
	rows, err := dbConnection.Query(`
		SELECT gr.request_id, gr.group_id, gr.user_id, gr.status, g.group_name
		FROM group_requests gr
		INNER JOIN groups g ON gr.group_id = g.group_id
		WHERE gr.user_id = $1
	`, userID)
	if err != nil {
		log.Printf("Error querying group requests sent by user from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and create group request objects with group names
	var groupRequestsWithName []GroupRequestWithName
	for rows.Next() {
		var groupRequest GroupRequestWithName
		if err := rows.Scan(&groupRequest.RequestID, &groupRequest.GroupID, &groupRequest.UserID, &groupRequest.Status, &groupRequest.GroupName); err != nil {
			log.Printf("Error scanning group request rows: %v", err)
			return nil, err
		}
		groupRequestsWithName = append(groupRequestsWithName, groupRequest)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over group request rows: %v", err)
		return nil, err
	}

	return groupRequestsWithName, nil
}
