package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/models"
)

// GetCommentsForPost retrieves all comments for a post with the given post ID
// GetCommentsForPost retrieves all comments for a post with the given post ID
func GetCommentsForPost(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Get the post ID from the URL parameters
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Missing post_id parameter", http.StatusBadRequest)
		return
	}

	// Get all comments for the specified post
	comments, err := GetCommentsByPostID(dbConnection, postID)
	if err != nil {
		// Log the error for debugging purposes
		fmt.Println("Error fetching comments:", err)
		// Return an HTTP response with a 500 Internal Server Error status
		http.Error(w, "Error getting comments for the post", http.StatusInternalServerError)
		return
	}





	// Respond with the comments in the JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}


// GetCommentsByPostID retrieves all comments for a post with the given post ID
// GetCommentsByPostID retrieves all comments for a post with the given post ID
func GetCommentsByPostID(db *sql.DB, postID string) ([]models.Comment, error) {
	rows, err := db.Query(`
		SELECT
			c.comment_id,
			c.content,
			c.comment_created_at,
			c.author_id,
			c.post_id,
			c.author_nickname,
			c.image
		FROM
			comments c
		WHERE
			c.post_id = ?;
	`, postID)
	if err != nil {
		return nil, fmt.Errorf("error fetching comments: %v", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.CommentID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.AuthorID,
			&comment.PostID,
			&comment.AuthorNickname,
			&comment.Image,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment rows: %v", err)
		}

		// Fetch likes count for the current comment
		likesCount, likeErr := GetLikesCountForComment(db, comment.CommentID)
		if likeErr != nil {
			return nil, fmt.Errorf("error fetching likes count for comment: %v", likeErr)
		}
		comment.LikesCount = likesCount

		comments = append(comments, comment)
	}

	return comments, nil
}


// GetLikesCountForComment retrieves the likes count for a specific comment
func GetLikesCountForComment(db *sql.DB, commentID string) (int, error) {
    fmt.Println("comment ID: ", commentID)
	var likesCount int
	err := db.QueryRow("SELECT COUNT(*) FROM CommentLikes WHERE comment_id = ?", commentID).Scan(&likesCount)
	if err != nil {
		return 0, fmt.Errorf("error fetching likes count for comment: %v", err)
	}

    fmt.Println("likes count: ", likesCount)
	return likesCount, nil
}
