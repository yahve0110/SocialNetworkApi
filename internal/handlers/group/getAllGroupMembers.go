package groupHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

type GroupMembersResponse struct {
	Members  []models.User
	IsMember bool
}

func GetAllGroupMembers(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Extract the group ID from the request parameters
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

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

	// Query the database to get all members of the group
	members, err := GetAllGroupMembersFromDatabase(dbConnection, groupID, userID)
	if err != nil {
		log.Printf("Error fetching group members from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the members into JSON and write it to the response writer
	if err := json.NewEncoder(w).Encode(members); err != nil {
		log.Printf("Error encoding group members to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetAllGroupMembersFromDatabase(dbConnection *sql.DB, groupID, userID string) (GroupMembersResponse, error) {
	var response GroupMembersResponse

	// Query the database to get all members of the group
	query := `
        SELECT u.user_id, u.username, u.profile_picture, u.first_name, u.last_name
        FROM users u
        JOIN group_members gm ON u.user_id = gm.user_id
        WHERE gm.group_id = $1`

	rows, err := dbConnection.Query(query, groupID)
	if err != nil {
		log.Printf("Error querying group members from database: %v", err)
		return response, err
	}
	defer rows.Close()

	// Iterate through the result set and create user objects
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.ProfilePicture, &user.FirstName, &user.LastName); err != nil {
			log.Printf("Error scanning user rows: %v", err)
			return response, err
		}
		response.Members = append(response.Members, user)
	}

	// Check if the current user is a member of the group
	for _, member := range response.Members {
		if member.UserID == userID {
			response.IsMember = true
			break
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over user rows: %v", err)
		return response, err
	}

	return response, nil
}
