package commentHandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/models"
)

// GetCommentsForPost retrieves all comments for a post with the given post ID
func GetCommentsForPost(w http.ResponseWriter, r *http.Request) {
	dbConnection := database.DB

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Missing post_id parameter", http.StatusBadRequest)
		return
	}

	comments, err := GetCommentsByPostID(dbConnection, postID)
	if err != nil {
		fmt.Println("Error fetching comments:", err)
		http.Error(w, "Error getting comments for the post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

// GetCommentsByPostID retrieves all comments for a post with the given post ID, sorted by creation date
func GetCommentsByPostID(db *sql.DB, postID string) ([]models.Comment, error) {
    query := `
    SELECT
        c.comment_id,
        c.content,
        c.comment_created_at,
        c.author_id,
        c.post_id,
        c.author_nickname,
        c.image,
        u.first_name,
        u.last_name,
        u.profile_picture
    FROM
        comments c
    JOIN
        users u ON c.author_id = u.user_id
    WHERE
        c.post_id = ?
    ORDER BY
        c.comment_created_at DESC;
    `

    rows, err := db.Query(query, postID)
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
            &comment.AuthorFirstName,
            &comment.AuthorLastName,
            &comment.AuthorAvatar,
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

    // Если нет комментариев, возвращаем пустой массив
    if len(comments) == 0 {
        return nil, nil
    }

    return comments, nil
}

// GetLikesCountForComment retrieves the likes count for a specific comment
func GetLikesCountForComment(db *sql.DB, commentID string) (int, error) {
	var likesCount int
	err := db.QueryRow("SELECT COUNT(*) FROM CommentLikes WHERE comment_id = ?", commentID).Scan(&likesCount)
	if err != nil {
		return 0, fmt.Errorf("error fetching likes count for comment: %v", err)
	}

	return likesCount, nil
}
