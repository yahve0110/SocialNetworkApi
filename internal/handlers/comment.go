package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
	"time"
	"github.com/google/uuid"
)



// Modify the AddComment function
func AddComment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("comment")

	var newComment models.Comment

	err := json.NewDecoder(r.Body).Decode(&newComment)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate content
	if newComment.Content == "" {
		http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	if newComment.PostID == "" {
		http.Error(w, "PostId cannot be empty", http.StatusBadRequest)
		return
	}

	// Generate a UUID for CommentID
	newComment.CommentID = uuid.New().String()

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Get the user ID and nickname based on the current user's session
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get the user ID and nickname based on the current user's session
	userID, err := helpers.GetUserIDFromSession(dbConnection, cookie.Value)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Set the author ID and nickname for the new comment
	newComment.AuthorID = userID

	// Fetch the author's nickname
	err = dbConnection.QueryRow("SELECT username FROM users WHERE user_id = ?", userID).Scan(&newComment.AuthorNickname)
	if err != nil {
		// Log the error for debugging purposes
		fmt.Println("Error fetching author's nickname:", err)
		// Return an HTTP response with a 500 Internal Server Error status
		http.Error(w, "Error fetching author's nickname", http.StatusInternalServerError)
		return
	}

	// Set the comment creation timestamp
	newComment.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	// Insert the new comment into the database
	_, err = dbConnection.Exec(`
		INSERT INTO comments (comment_id, content, comment_created_at, author_id, post_id, author_nickname, image)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, newComment.CommentID, newComment.Content, newComment.CreatedAt, newComment.AuthorID, newComment.PostID, newComment.AuthorNickname, newComment.Image)
	if err != nil {
		// Log the error for debugging purposes
		fmt.Println("Error inserting comment into database:", err)
		// Return an HTTP response with a 500 Internal Server Error status
		http.Error(w, "Error inserting comment into database", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Comment created with ID %s", newComment.CommentID)))
}