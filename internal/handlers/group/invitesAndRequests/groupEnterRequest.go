package groupInviteHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"

	"github.com/google/uuid"
)

// SendGroupRequestHandler handles the sending of group join requests

func SendGroupRequestHandler(w http.ResponseWriter, r *http.Request) {
	var requestData models.GroupRequest

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Check if the requester is authenticated
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

	// Check if the user is already a member of the group
	isMember, err := helpers.IsUserGroupMember(dbConnection, userID, requestData.GroupID)
	if err != nil {
		log.Printf("Error checking if user is already a group member: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if isMember {
		http.Error(w, "User is already a member of the group", http.StatusBadRequest)
		return
	}

	// Check if the user has already sent a request to join the group
	requestExists, err := helpers.GroupRequestExists(dbConnection, requestData.GroupID, userID)
	if err != nil {
		log.Printf("Error checking if user has already sent a request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if requestExists {
		http.Error(w, "User has already sent a request to join the group", http.StatusBadRequest)
		return
	}

	// Set the user ID for the request
	requestData.UserID = userID

	// Set the status for the request
	requestData.Status = "Pending"

	// Generate a UUID for the request ID
	requestData.RequestID = uuid.New().String()

	// Save the request to the database
	err = saveGroupRequestToDatabase(requestData, dbConnection)
	if err != nil {
		log.Printf("Error saving group request to database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(requestData); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// saveGroupRequestToDatabase saves a group request to the database
func saveGroupRequestToDatabase(requestData models.GroupRequest, db *sql.DB) error {

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO group_requests (request_id, group_id, user_id, status) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(requestData.RequestID, requestData.GroupID, requestData.UserID, requestData.Status)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}
