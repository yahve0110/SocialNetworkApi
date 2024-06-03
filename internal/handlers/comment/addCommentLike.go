package commentHandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

type CommentLike struct {
	CommentID string `json:"comment_id"`
}

func AddCommentLike(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Get user ID from session (assuming session ID is stored in a cookie)
	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "Failed to get session ID from cookie", http.StatusBadRequest)
		return
	}

	userID, err := helpers.GetUserIDFromSession(dbConnection, sessionID.Value)
	if err != nil {
		http.Error(w, "Failed to get user ID from session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the request body to get comment-like data
	var commentLike CommentLike
	err = json.NewDecoder(r.Body).Decode(&commentLike)
	if err != nil {
		http.Error(w, "Failed to decode request body Comment ID cannot be empty: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user has already liked the comment
	var existingUserID string
	err = dbConnection.QueryRow("SELECT user_id FROM CommentLikes WHERE comment_id  = ? AND user_id = ?", commentLike.CommentID, userID).Scan(&existingUserID)
	if err == nil {
		_, err = dbConnection.Exec("DELETE FROM CommentLikes WHERE comment_id = ? AND user_id = ?", commentLike.CommentID, userID)
	} else {
		_, err = dbConnection.Exec("INSERT INTO CommentLikes (comment_id, user_id) VALUES (?, ?)", commentLike.CommentID, userID)
	}

	if err != nil {
		http.Error(w, "Cannot set like: "+err.Error(), http.StatusInternalServerError)
		return
	}

	likesCount, err := GetCommentLikesCount(dbConnection, commentLike.CommentID)
	if err != nil {
		http.Error(w, "Failed to get comment likes count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Success    bool `json:"success"`
		LikesCount int  `json:"likes_count"`
	}{
		Success:    true,
		LikesCount: likesCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetCommentLikesCount(dbConnection *sql.DB, commentID string) (int, error) {
	var likesCount int
	err := dbConnection.QueryRow("SELECT COUNT(*) FROM CommentLikes WHERE comment_id = ?", commentID).Scan(&likesCount)
	if err != nil {
		return 0, fmt.Errorf("error fetching comment likes count: %v", err)
	}
	return likesCount, nil
}
