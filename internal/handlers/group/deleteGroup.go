package groupHandlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

// DeleteGroupHandler handles the deletion of a group
func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the database package
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



	// Check if the user is the creator of the group
	groupID := r.FormValue("group_id")
	log.Printf("IsUserGroupCreator - Group ID: %s", groupID)
	isGroupCreator, err := IsUserGroupCreator(dbConnection, userID, groupID)
	if err != nil {
		log.Printf("Error checking group creator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("IsUserGroupCreator - exists: %v", isGroupCreator)

	if !isGroupCreator {
		http.Error(w, "Unauthorized: Only group creator can delete the group", http.StatusUnauthorized)
		return
	}

	// Delete the group and associated data
	err = DeleteGroupAndAssociatedData(dbConnection, groupID)
	if err != nil {
		log.Printf("Error deleting group and associated data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Group deleted successfully"))
}

// IsUserGroupCreator checks if a user is the creator of the specified group
func IsUserGroupCreator(db *sql.DB, userID, groupID string) (bool, error) {
	// Query the groups table to check if the user is the creator
	query := "SELECT creator_id FROM groups WHERE group_id = ?"
	var groupCreatorID string
	err := db.QueryRow(query, groupID).Scan(&groupCreatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("group with ID %s not found", groupID)
		}
		return false, err
	}

	// Check if the group creator ID matches the user ID
	isCreator := groupCreatorID == userID
	return isCreator, nil
}

// DeleteGroupAndAssociatedData deletes a group and all associated data from the database, including images from the cloud
func DeleteGroupAndAssociatedData(db *sql.DB, groupID string) error {
	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete all posts related to the group
	if err := deletePostsByGroupID(tx, groupID, db); err != nil {
		tx.Rollback()
		return err
	}

	// Delete all events related to the group
	if err := deleteEventsByGroupID(tx, groupID, db); err != nil {
		tx.Rollback()
		return err
	}

	// Delete the group itself
	if _, err := tx.Exec("DELETE FROM groups WHERE group_id = ?", groupID); err != nil {
		tx.Rollback()
		return err
	}

	// Delete the group's image from cloud if it exists
	groupImageURL, err := GetGroupImageURL(db, groupID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if groupImageURL != "" && groupImageURL != "https://res.cloudinary.com/djkotlye3/image/upload/v1713162945/g0n2phibtawxxgwmxnig.jpg" {
		if err := helpers.DeleteFromCloudinary(groupImageURL); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func deleteImagesFromCloud(imageURLs []string) error {
	for _, imageURL := range imageURLs {
		if imageURL != "" {
			if err := helpers.DeleteFromCloudinary(imageURL); err != nil {
				return err
			}
		}
	}
	return nil
}

func deletePostsByGroupID(tx *sql.Tx, groupID string, db *sql.DB) error {
	rows, err := tx.Query("SELECT post_id FROM group_posts WHERE group_id = ?", groupID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var postID string
		if err := rows.Scan(&postID); err != nil {
			return err
		}
		if err := deleteCommentsByPostID(tx, postID, db); err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec("DELETE FROM group_posts WHERE group_id = ?", groupID)
	if err != nil {
		return err
	}

	return nil
}

func deleteEventsByGroupID(tx *sql.Tx, groupID string, db *sql.DB) error {
	rows, err := tx.Query("SELECT event_img FROM group_events WHERE group_id = ?", groupID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var imageURLs []string
	for rows.Next() {
		var imageURL string
		if err := rows.Scan(&imageURL); err != nil {
			return err
		}
		imageURLs = append(imageURLs, imageURL)
	}

	if err := deleteImagesFromCloud(imageURLs); err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM group_events WHERE group_id = ?", groupID)
	if err != nil {
		return err
	}

	return nil
}

func deleteCommentsByPostID(tx *sql.Tx, postID string, db *sql.DB) error {
	rows, err := tx.Query("SELECT image FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var imageURLs []string
	for rows.Next() {
		var imageURL string
		if err := rows.Scan(&imageURL); err != nil {
			return err
		}
		imageURLs = append(imageURLs, imageURL)
	}

	if err := deleteImagesFromCloud(imageURLs); err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return err
	}

	return nil
}

// GetGroupImageURL retrieves the image URL of a group from the database
func GetGroupImageURL(db *sql.DB, groupID string) (string, error) {
	// Prepare SQL statement to select the image URL of the group
	stmt, err := db.Prepare("SELECT group_image FROM groups WHERE group_id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	// Execute the SQL statement
	var imageURL string
	err = stmt.QueryRow(groupID).Scan(&imageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			// Group with the specified ID not found
			return "", nil
		}
		return "", err
	}

	return imageURL, nil
}
