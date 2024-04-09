package groupEventHandlers

import (
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

// NotJoinGroupEventHandler handles a user not joining (not going) an event
func DeclineEventHandler(w http.ResponseWriter, r *http.Request) {
    var notJoinRequest EventJoinRequest

    err := json.NewDecoder(r.Body).Decode(&notJoinRequest)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

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

    // Check if the event exists
    exists, err := IsEventExist(dbConnection, notJoinRequest.EventID)
    if err != nil {
        log.Printf("Error checking if event exists: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    if !exists {
        http.Error(w, "Event not found", http.StatusNotFound)
        return
    }

    // Check if the user has already not joined the event
    alreadyNotJoined, err := IsUserNotJoinedEvent(dbConnection, userID, notJoinRequest.EventID)
    if err != nil {
        log.Printf("Error checking if user has already declined the event: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    if alreadyNotJoined {
        http.Error(w, "User has already declined the event", http.StatusConflict)
        return
    }

    // Check if the user has already joined the event
    alreadyJoined, err := IsUserJoinedEvent(dbConnection, userID, notJoinRequest.EventID)
    if err != nil {
        log.Printf("Error checking if user has already joined the event: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    if alreadyJoined {
        http.Error(w, "User has already joined the event", http.StatusConflict)
        return
    }

    // Add the user to the not going members of the event
    err = AddUserToNotGoingMembers(dbConnection, userID, notJoinRequest.EventID)
    if err != nil {
        log.Printf("Error adding user to not going members: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User declined event"))
}
