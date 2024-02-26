package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Content string `json:"content"`
	Image   string `json:"image"`
}

func AddPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("post")

	var newPost models.Post

	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate content
	if newPost.Content == "" {
		http.Error(w, "Post content cannot be empty", http.StatusBadRequest)
		return
	}
	// Generate a UUID for PostID
	newPost.PostID = uuid.New().String()

	//set post privacy by default public
	if newPost.Private == "" {
		newPost.Private = "public"
	}

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

	// Set the author ID for the new post
	newPost.AuthorID = userID

	// Set the post creation timestamp
	newPost.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	// Call the function to insert the post into the database
	if err := InsertPost(dbConnection, newPost); err != nil {
		// Log the error for debugging purposes
		fmt.Println("Error inserting post into database:", err)
		// Return an HTTP response with a 500 Internal Server Error status
		http.Error(w, "Error inserting post into database", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Post created with ID %s", newPost.PostID)))
}

// InsertPost inserts a new post into the database
func InsertPost(db *sql.DB, post models.Post) error {
	_, err := db.Exec(`
		INSERT INTO posts (post_id, author_id, content, post_created_at, likes_count,privacy, image)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, post.PostID, post.AuthorID, post.Content, post.CreatedAt, post.LikesCount,post.Private, post.Image)
	if err != nil {
		return fmt.Errorf("error inserting post into database: %v", err)
	}
	return nil
}
