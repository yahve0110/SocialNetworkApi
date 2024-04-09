package groupPostCommentHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/models"
)

// GetGroupPostCommentsHandler handles the retrieval of comments for a group post
func GetGroupPostCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	postID := r.FormValue("post_id")

	// Retrieve all comments and likes for the group post
	comments, err := getGroupPostCommentsFromDatabase(dbConnection, postID)
	if err != nil {
		log.Printf("Error fetching comments from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// getGroupPostCommentsFromDatabase fetches all comments and likes for a group post from the database
func getGroupPostCommentsFromDatabase(dbConnection *sql.DB, postID string) ([]models.GroupPostComment, error) {
	// Query all comments and like counts for the group post from the "group_comments" table
	rows, err := dbConnection.Query(`
		SELECT
			c.comment_id,
			c.post_id,
			c.author_id,
			c.content,
			c.comment_date,
			(SELECT COUNT(*) FROM group_comment_likes l WHERE l.comment_id = c.comment_id) AS like_count
		FROM group_comments c
		WHERE c.post_id = ?
	`, postID)
	if err != nil {
		log.Printf("Error querying comments from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and create comment objects
	var comments []models.GroupPostComment
	for rows.Next() {
		var comment models.GroupPostComment
		if err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt, &comment.LikeCount); err != nil {
			log.Printf("Error scanning comment rows: %v", err)
			return nil, err
		}
		comments = append(comments, comment)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over comment rows: %v", err)
		return nil, err
	}

	return comments, nil
}
