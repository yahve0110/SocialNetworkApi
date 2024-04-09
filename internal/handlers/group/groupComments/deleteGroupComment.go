package groupPostCommentHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

// DeleteGroupCommentHandler handles the deletion of a group comment
func DeleteGroupCommentHandler(w http.ResponseWriter, r *http.Request) {
	var requestData models.GroupPostComment

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Check if the requester is authenticated
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

	// Check if the user is the creator of the comment or the group creator
	isCommentCreator, err := IsUserCommentCreator(dbConnection, userID, requestData.CommentID)
	if err != nil {
		log.Printf("Error checking comment creator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	isGroupCreator, err := helpers.IsUserGroupCreator(dbConnection, userID, requestData.GroupID)
	if err != nil {
		log.Printf("Error checking group creator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isCommentCreator && !isGroupCreator {
		http.Error(w, "Unauthorized: Only comment creator or group creator can delete the comment", http.StatusUnauthorized)
		return
	}

	// Delete the group comment
	err = DeleteGroupComment(dbConnection, requestData.CommentID)
	if err != nil {
		log.Printf("Error deleting group comment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestData)
}

// IsUserCommentCreator checks if a user is the creator of the specified comment
func IsUserCommentCreator(db *sql.DB, userID, commentID string) (bool, error) {
	// Query the group_comments table to check if the user is the creator
	query := "SELECT EXISTS(SELECT 1 FROM group_comments WHERE author_id = ? AND comment_id = ?)"
	var exists bool
	err := db.QueryRow(query, userID, commentID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking comment creator: %v", err)
		return false, err
	}

	return exists, nil
}

// DeleteGroupComment deletes a comment from the database
func DeleteGroupComment(db *sql.DB, commentID string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("DELETE FROM group_comments WHERE comment_id = ?")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(commentID)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}
