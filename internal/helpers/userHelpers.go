package helpers

import (
	"database/sql"
	"fmt"
)

// GetUserIDFromSession retrieves the user ID from the database based on the session ID
func GetUserIDFromSession(db *sql.DB, sessionID string) (string, error) {
	var userID string

	// Query the database to get the user ID associated with the session ID
	err := db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ? AND expiration_time > CURRENT_TIMESTAMP", sessionID).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("error getting user ID for session ID %s: %v", sessionID, err)
	}

	return userID, nil
}
