package helpers

import (
	"database/sql"
	"fmt"
	"social/internal/models"
)

// GetFollowers retrieves all followers for a given user ID
func GetFollowers(db *sql.DB, userID string) ([]models.User, error) {
	rows, err := db.Query(`
		SELECT
			users.user_id,
			users.username,
			users.first_name,
			users.last_name,
			users.email,
			users.gender,
			users.birth_date,
			users.profile_picture,
			users.about
		FROM
			Followers
		JOIN users ON Followers.user_following = users.user_id
		WHERE
			Followers.user_followed = ? AND Followers.user_followed_status = 'accepted'
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching followers: %v", err)
	}
	defer rows.Close()

	var followers []models.User
	for rows.Next() {
		var follower models.User
		err := rows.Scan(
			&follower.UserID,
			&follower.Username,
			&follower.FirstName,
			&follower.LastName,
			&follower.Email,
			&follower.Gender,
			&follower.BirthDate,
			&follower.ProfilePicture,
			&follower.About,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning follower rows: %v", err)
		}
		followers = append(followers, follower)
	}

	return followers, nil
}

// GetFollowing retrieves all users a given user is following
func GetFollowing(db *sql.DB, userID string) ([]models.User, error) {
	rows, err := db.Query(`
		SELECT
			users.user_id,
			users.username,
			users.first_name,
			users.last_name,
			users.email,
			users.gender,
			users.birth_date,
			users.profile_picture,
			users.about
		FROM
			Followers
		JOIN users ON Followers.user_followed = users.user_id
		WHERE
			Followers.user_following = ? AND Followers.user_followed_status = 'accepted'
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching following users: %v", err)
	}
	defer rows.Close()

	var following []models.User
	for rows.Next() {
		var followedUser models.User
		err := rows.Scan(
			&followedUser.UserID,
			&followedUser.Username,
			&followedUser.FirstName,
			&followedUser.LastName,
			&followedUser.Email,
			&followedUser.Gender,
			&followedUser.BirthDate,
			&followedUser.ProfilePicture,
			&followedUser.About,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning following rows: %v", err)
		}
		following = append(following, followedUser)
	}

	return following, nil
}

// DeleteFollowStatus deletes a follow status in the Followers table
func DeleteFollowStatus(db *sql.DB, followerID, followedID string) error {
	_, err := db.Exec("DELETE FROM Followers WHERE user_following = ? AND user_followed = ?", followerID, followedID)
	if err != nil {
		return fmt.Errorf("error deleting follow status: %v", err)
	}
	return nil
}

// GetFollowersWithStatus retrieves all followers with a specific status from the Followers table
func GetFollowersWithStatus(db *sql.DB, userFollowed, status string) ([]string, error) {
	rows, err := db.Query("SELECT user_following FROM Followers WHERE user_followed = ? AND user_followed_status = ?", userFollowed, status)
	if err != nil {
		return nil, fmt.Errorf("error fetching followers with status: %v", err)
	}
	defer rows.Close()

	var followers []string
	for rows.Next() {
		var followerID string
		err := rows.Scan(&followerID)
		if err != nil {
			return nil, fmt.Errorf("error scanning follower rows: %v", err)
		}
		followers = append(followers, followerID)
	}

	return followers, nil
}

func AcceptPendingFollower(db *sql.DB, userFollowingID, userFollowedID string) error {
	// Check if the pending follower request exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Followers WHERE user_followed = ? AND user_following = ? AND user_followed_status = 'pending'", userFollowedID, userFollowingID).Scan(&count)
	if err != nil {
		fmt.Printf("Error checking pending follower request existence: %v\n", err)
		return fmt.Errorf("error checking pending follower request existence: %v", err)
	}

	if count == 0 {
		// The pending follower request doesn't exist
		fmt.Printf("Pending follower request not found for userFollowedID: %s, userFollowingID: %s\n", userFollowedID, userFollowingID)
		return fmt.Errorf("pending follower request not found for userFollowedID: %s, userFollowingID: %s", userFollowedID, userFollowingID)
	}

	// Update the status to 'accepted'
	query := "UPDATE Followers SET user_followed_status = 'accepted' WHERE user_followed = ? AND user_following = ?"
	_, err = db.Exec(query, userFollowedID, userFollowingID)
	if err != nil {
		fmt.Printf("Error executing SQL query: %s with userFollowedID: %s, userFollowingID: %s\n", query, userFollowedID, userFollowingID)
		fmt.Printf("Error updating pending follower request: %v\n", err)
		return fmt.Errorf("error updating pending follower request: %v", err)
	}

	fmt.Printf("Pending follower request accepted: userFollowedID %s, userFollowingID %s\n", userFollowedID, userFollowingID)

	return nil
}




