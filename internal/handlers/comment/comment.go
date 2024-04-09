package commentHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
	"time"
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

	// Set the author ID for the new comment
	newComment.AuthorID = userID

	// Fetch the author's first name and last name
	var authorFirstName, authorLastName, authorAvatar string
	err = dbConnection.QueryRow("SELECT first_name, last_name, profile_picture FROM users WHERE user_id = ?", userID).Scan(&authorFirstName, &authorLastName, &authorAvatar)
	if err != nil {
		// Log the error for debugging purposes
		fmt.Println("Error fetching author's information:", err)
		// Return an HTTP response with a 500 Internal Server Error status
		http.Error(w, "Error fetching author's information", http.StatusInternalServerError)
		return
	}


    //upload comment image to cloud storage
    commentImageBase64 := newComment.Image
	if commentImageBase64 != "" {
        cloudinaryURL, err := helpers.ImageToCloud(commentImageBase64, w)
        if err != nil {
            // Handle error
            return
        }
        newComment.Image = cloudinaryURL
    }


	// Set the author's first name and last name for the new comment
	newComment.AuthorFirstName = authorFirstName
	newComment.AuthorLastName = authorLastName
	newComment.AuthorAvatar = authorAvatar

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

	// Respond with the created comment in the response body
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newComment)
}
