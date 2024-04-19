package groupPostHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

// LikeGroupPostHandler handles the liking of a group post
func LikeGroupPostHandler(w http.ResponseWriter, r *http.Request) {
	var requestData models.GroupPostLike

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the db package
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

	// Check if the user has already liked the post
	liked, err := HasUserLikedPost(dbConnection, userID, requestData.PostID)
	if err != nil {
		log.Printf("Error checking if user has liked the post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// If already liked, remove the like; otherwise, add a like
	if liked {
		// Remove the like from the group post
		err = RemoveLikeFromGroupPost(dbConnection, userID, requestData.PostID)
		if err != nil {
			log.Printf("Error removing like from group post: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		// Add a like to the group post
		err = AddLikeToGroupPost(dbConnection, userID, requestData.PostID)
		if err != nil {
			log.Printf("Error adding like to group post: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Get the count of likes for the post
	likesCount, err := GetLikesCountForPost(dbConnection, requestData.PostID)
	if err != nil {
		log.Printf("Error fetching like count for post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return the count of likes in the response
	responseData := map[string]interface{}{
		"likesCount": likesCount,
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

// HasUserLikedPost checks if a user has already liked a group post
func HasUserLikedPost(db *sql.DB, userID, postID string) (bool, error) {
	// Query the group_post_likes table to check if the user has already liked the post
	query := "SELECT EXISTS(SELECT 1 FROM group_post_likes WHERE user_id = ? AND post_id = ?)"
	var exists bool
	err := db.QueryRow(query, userID, postID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if user has liked the post: %v", err)
		return false, err
	}

	return exists, nil
}

// AddLikeToGroupPost adds a like to a group post
func AddLikeToGroupPost(db *sql.DB, userID, postID string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO group_post_likes (post_id, user_id) VALUES (?, ?)")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(postID, userID)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}

// RemoveLikeFromGroupPost removes a like from a group post
func RemoveLikeFromGroupPost(db *sql.DB, userID, postID string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("DELETE FROM group_post_likes WHERE user_id = ? AND post_id = ?")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(userID, postID)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}


// GetLikesCountForPost fetches the count of likes for a group post
func GetLikesCountForPost(db *sql.DB, postID string) (int, error) {
	// Query the group_post_likes table to get the count of likes for the post
	query := "SELECT COUNT(*) FROM group_post_likes WHERE post_id = ?"
	var likesCount int
	err := db.QueryRow(query, postID).Scan(&likesCount)
	if err != nil {
		log.Printf("Error fetching like count for post %s: %v", postID, err)
		return 0, err
	}

	return likesCount, nil
}
