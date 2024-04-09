package groupPostCommentHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

// LikeGroupCommentHandler handles the liking or unliking of a group comment
func LikeGroupCommentHandler(w http.ResponseWriter, r *http.Request) {
	var requestData models.GroupCommentLike

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

	// Check if the user has already liked the comment
	liked, err := HasUserLikedComment(dbConnection, userID, requestData.CommentID)
	if err != nil {
		log.Printf("Error checking if user has liked the comment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// If already liked, remove the like; otherwise, add a like
	if liked {
		// Remove the like from the group comment
		err = RemoveLikeFromGroupComment(dbConnection, userID, requestData.CommentID)
		if err != nil {
			log.Printf("Error removing like from group comment: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		// Add a like to the group comment
		err = AddLikeToGroupComment(dbConnection, userID, requestData.CommentID)
		if err != nil {
			log.Printf("Error adding like to group comment: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestData)
}

// HasUserLikedComment checks if a user has already liked a group comment
func HasUserLikedComment(db *sql.DB, userID, commentID string) (bool, error) {
	// Query the group_comment_likes table to check if the user has already liked the comment
	query := "SELECT EXISTS(SELECT 1 FROM group_comment_likes WHERE user_id = ? AND comment_id = ?)"
	var exists bool
	err := db.QueryRow(query, userID, commentID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if user has liked the comment: %v", err)
		return false, err
	}

	return exists, nil
}

// AddLikeToGroupComment adds a like to a group comment
func AddLikeToGroupComment(db *sql.DB, userID, commentID string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO group_comment_likes (comment_id, user_id) VALUES (?, ?)")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()


	// Execute the SQL statement
	_, err = stmt.Exec(commentID, userID)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}

// RemoveLikeFromGroupComment removes a like from a group comment
func RemoveLikeFromGroupComment(db *sql.DB, userID, commentID string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("DELETE FROM group_comment_likes WHERE user_id = ? AND comment_id = ?")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(userID, commentID)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}
