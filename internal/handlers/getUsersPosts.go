package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

func GetUserPosts(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
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

	// Get all posts created by the user
	userPosts, err := GetPostsByUserID(dbConnection, userID)
	if err != nil {
		// Log the error for debugging purposes
		fmt.Println("error fetching posts:", err)
		// Return an HTTP response with a 500 Internal Server Error status
		http.Error(w, "Error getting user posts", http.StatusInternalServerError)
		return
	}

	// Respond with the user's posts in the JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userPosts)
}

// GetPostsByUserID retrieves all posts created by a user with the given user ID
func GetPostsByUserID(db *sql.DB, userID string) ([]models.Post, error) {
	rows, err := db.Query(`
	SELECT
	posts.post_id,
    users.user_id,
    users.username,
    posts.content,
    posts.post_created_at,
	COUNT(postLikes.user_id) AS tlikes_count,
    posts.image
FROM
    posts
JOIN users ON posts.author_id = users.user_id
LEFT JOIN postLikes ON posts.post_id = postLikes.post_id
WHERE
    posts.author_id = ?
GROUP BY
    posts.post_id, users.user_id, users.username, posts.content, posts.post_created_at, posts.likes_count, posts.image
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching posts: %v", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.PostID,
			&post.AuthorID,
			&post.AuthorNickname,
			&post.Content,
			&post.CreatedAt,
			&post.LikesCount, // Use LikesCount for total_likes
			&post.Image,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post rows: %v", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}
