package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"time"
)

type SessionRequest struct {
	SessionID string `json:"sessionId"`
}

func IsSessionValid(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body into SessionRequest struct
	var sessionReq SessionRequest
	err := json.NewDecoder(r.Body).Decode(&sessionReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sessionID := sessionReq.SessionID

	dbConnection := database.DB

	// Query the sessions table to check if the session is valid
	var expirationTime time.Time
	err = dbConnection.QueryRow("SELECT expiration_time FROM sessions WHERE session_id = ?", sessionID).Scan(&expirationTime)

	if err != nil {
		// Session ID not found or other database error
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if the session is expired
	if time.Now().After(expirationTime) {
		// Session has expired
		http.Error(w, "Session Expired", http.StatusUnauthorized)
		return
	}

	// Session is valid
	fmt.Fprintf(w, "Session is valid")
}
