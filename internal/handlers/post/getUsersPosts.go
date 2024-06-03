package postHandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/models"
)

func GetUserPosts(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Parse the request body to get the user ID
	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, exists := requestBody["user_id"]
	if !exists || userID == "" {
		http.Error(w, "User ID is required in the request body", http.StatusBadRequest)
		return
	}



	// Get all posts created by the specified user
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

// GetPostsByUserID retrieves all posts created by a user with the given user ID, sorted by creation date (newest first)
func GetPostsByUserID(db *sql.DB, userID string) ([]models.Post, error) {
    query := `
        SELECT
            posts.post_id,
            users.user_id,
            users.username,
            users.first_name,
            users.last_name,
            posts.content,
            posts.post_created_at,
            posts.privacy,
            COUNT(postLikes.user_id) AS tlikes_count,
            posts.image
        FROM
            posts
        JOIN users ON posts.author_id = users.user_id
        LEFT JOIN postLikes ON posts.post_id = postLikes.post_id
        WHERE
            posts.author_id = ?
        GROUP BY
            posts.post_id, users.user_id, users.username, users.first_name, users.last_name, posts.content, posts.post_created_at, posts.likes_count, posts.image
        ORDER BY
            posts.post_created_at DESC
    `
    rows, err := db.Query(query, userID)
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
            &post.AuthorFirstName,
            &post.AuthorLastName,
            &post.Content,
            &post.CreatedAt,
            &post.Private,
            &post.LikesCount,
            &post.Image,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning post rows: %v", err)
        }

        if post.Private == "almost private" {
            allowedUsers, err := GetAllowedUsersForPost(db, post.PostID)
            if err != nil {
                return nil, fmt.Errorf("error getting allowed users for post ID %s: %v", post.PostID, err)
            }
            post.PrivateUsersArr = allowedUsers
        }

        // Add the post to the list
        posts = append(posts, post)
    }

    // Check if posts were obtained
    if len(posts) == 0 {
        return nil, fmt.Errorf("no posts found for user ID: %s", userID)
    }

    return posts, nil
}

func GetAllowedUsersForPost(db *sql.DB, postID string) ([]string, error) {
    query := `SELECT user_id FROM post_permissions WHERE post_id = ?`
    rows, err := db.Query(query, postID)
    if err != nil {
        return nil, fmt.Errorf("error fetching allowed users for post ID %s: %v", postID, err)
    }
    defer rows.Close()

    var allowedUsers []string
    for rows.Next() {
        var userID string
        err := rows.Scan(&userID)
        if err != nil {
            return nil, fmt.Errorf("error scanning allowed users rows for post ID %s: %v", postID, err)
        }
        allowedUsers = append(allowedUsers, userID)
    }

    return allowedUsers, nil
}
