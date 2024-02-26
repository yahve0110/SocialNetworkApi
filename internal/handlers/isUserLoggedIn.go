package handlers

import (
	"fmt"
	"net/http"
	database "social/internal/db"
	"time"
)

func IsSessionValid(w http.ResponseWriter, r *http.Request) {
	// Retrieve session ID from cookies
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		// Cookie not found, session is invalid
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID := cookie.Value

	fmt.Println("Session ID: ", sessionID)
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
