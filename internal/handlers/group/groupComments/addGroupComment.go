package groupPostCommentHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"

	"github.com/google/uuid"
)

func AddGroupPostCommentHandler(w http.ResponseWriter, r *http.Request) {
	var requestData models.GroupPostComment

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Check if the requester is a group member or the post creator
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

	// Check if the user is a group member or the post creator
	isGroupMember, err := helpers.IsUserGroupMember(dbConnection, userID, requestData.GroupID)
	if err != nil {
		log.Printf("Error checking group membership: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	isPostCreator, err := helpers.IsUserPostCreator(dbConnection, userID, requestData.PostID)
	if err != nil {
		log.Printf("Error checking post creator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isGroupMember && !isPostCreator {
		http.Error(w, "Unauthorized: Only group members or post creator can leave comments", http.StatusUnauthorized)
		return
	}



	// Add the comment to the group post
	err = AddGroupPostComment(dbConnection, requestData.PostID, userID, requestData.Content)
	if err != nil {
		log.Printf("Error adding comment to group post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestData)
}

// AddGroupPostComment inserts a comment into the group_comments table
func AddGroupPostComment(db *sql.DB, postID, authorID, content string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO group_comments(comment_id, post_id, author_id, content, comment_date) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	commentID := uuid.New().String()

	// Execute the SQL statement
	_, err = stmt.Exec(commentID, postID, authorID, content)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}
