package helpers

import (
	"database/sql"
	"fmt"
	"log"
	"social/internal/models"
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

// UserExists checks if a user with the given user ID exists in the database
func UserExists(db *sql.DB, userID string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE user_id = ?)", userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %v", err)
	}
	return exists, nil
}

// IsUserProfilePublic checks if the user's profile is public based on the privacy setting
func IsUserProfilePublic(db *sql.DB, userID string) (string, error) {
	var privacy string
	err := db.QueryRow("SELECT privacy FROM users WHERE user_id = ?", userID).Scan(&privacy)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where the user is not found
			return "private", nil // Assuming private if user not found
		}
		return "", fmt.Errorf("error checking user profile privacy: %v", err)
	}
	return privacy, nil
}

func InsertFollowRequest(db *sql.DB, followerID, followedID string) error {
	_, err := db.Exec("INSERT INTO Followers (user_followed_status, user_followed,user_following) VALUES (?, ?, ?)", "pending", followedID, followerID)
	if err != nil {
		return fmt.Errorf("error inserting follow request: %v", err)
	}
	return nil
}

// InsertOrUpdateFollowStatus inserts or updates a follow status in the Followers table
func InsertOrUpdateFollowStatus(db *sql.DB, followerID, followedID string) error {
	// Check if the follow status already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Followers WHERE user_followed = ? AND user_following = ?", followedID, followerID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking follow status existence: %v", err)
	}

	if count > 0 {
		// Follow status exists, update it
		_, err := db.Exec("UPDATE Followers SET user_followed_status = 'accepted' WHERE user_followed = ? AND user_following = ?", followedID, followerID)
		if err != nil {
			log.Printf("Error updating follow status: %v", err)

			return fmt.Errorf("error updating follow status: %v", err)
		}
	} else {
		/// Follow status doesn't exist, insert it
		_, err := db.Exec("INSERT INTO Followers (user_followed_status, user_followed, user_following) VALUES (?, ?, ?)", "accepted", followedID, followerID)
		if err != nil {
			log.Printf("Error inserting follow status: %v", err)
			return fmt.Errorf("error inserting follow status: %v", err)
		}
	}

	return nil
}

// GetUserByID retrieves user information by user ID
func GetUserByID(db *sql.DB, userID string) (models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT user_id, username, first_name, last_name, email, gender, birth_date, profile_picture, about FROM users WHERE user_id = ?", userID).
		Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Gender, &user.BirthDate, &user.ProfilePicture, &user.About)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetUserIDByUsername(db *sql.DB, username string) (string, error) {
    var userID string

    // Prepare the SQL statement to select the user ID based on the username
    query := "SELECT user_id FROM users WHERE username = ? LIMIT 1"
    err := db.QueryRow(query, username).Scan(&userID)
    if err != nil {
        return "", err
    }

    return userID, nil
}