package groupEventHandlers

import (
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

type EventJoinRequest struct {
	EventID  string `json:"event_id"`
	MemberID string `json:"member_id"`
}

// JoinGroupEventHandler handles a user joining (going) an event
func JoinGroupEventHandler(w http.ResponseWriter, r *http.Request) {
    var joinRequest EventJoinRequest

    err := json.NewDecoder(r.Body).Decode(&joinRequest)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Access the global database connection from the database package
    dbConnection := database.DB

    // Check if the event exists
    exists, err := IsEventExist(dbConnection, joinRequest.EventID)
    if err != nil {
        log.Printf("Error checking if event exists: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    if !exists {
        http.Error(w, "Event not found", http.StatusNotFound)
        return
    }

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

    // Check if the user has already joined or not joined the event
    alreadyJoined, err := IsUserJoinedEvent(dbConnection, userID, joinRequest.EventID)
    if err != nil {
        log.Printf("Error checking if user has already joined the event: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    if alreadyJoined {
        http.Error(w, "User has already joined the event", http.StatusConflict)
        return
    }

    alreadyNotJoined, err := IsUserNotJoinedEvent(dbConnection, userID, joinRequest.EventID)
    if err != nil {
        log.Printf("Error checking if user has already declined event: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    if alreadyNotJoined {
        http.Error(w, "User has already declined event", http.StatusConflict)
        return
    }

    // Add the user to the going members of the event
    err = AddUserToGoingMembers(dbConnection, userID, joinRequest.EventID)
    if err != nil {
        log.Printf("Error adding user to going members: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User joined the event"))
}
