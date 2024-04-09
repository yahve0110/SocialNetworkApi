package groupPostHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"

	"github.com/google/uuid"
)

// CreateGroupPostHandler handles the creation of posts in a group
func CreateGroupPostHandler(w http.ResponseWriter, r *http.Request) {
	var postData models.GroupPost

	// Decode the request body into postData
	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Check if the user is authenticated
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

	// Check if the user is a member or the creator of the group
	isGroupMember, err := helpers.IsUserGroupMember(dbConnection, userID, postData.GroupID)
	if err != nil {
		log.Printf("Error checking group membership: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	isGroupCreator, err := helpers.IsUserGroupCreator(dbConnection, userID, postData.GroupID)
	if err != nil {
		log.Printf("Error checking group creator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isGroupMember && !isGroupCreator {
		http.Error(w, "Unauthorized: Only group members or creator can create posts", http.StatusUnauthorized)
		return
	}

	// Set the UserID and creation time for the post
	postData.AuthorID = userID
	postData.CreatedAt = time.Now()

	//create postId
	postData.PostID = uuid.New().String()

	// Insert the post into the database
	err = InsertGroupPost(dbConnection, postData)
	if err != nil {
		log.Printf("Error inserting group post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postData)
}

// InsertGroupPost inserts a new post into the database
func InsertGroupPost(db *sql.DB, post models.GroupPost) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO group_posts (post_id, group_id, author_id, content, post_date) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(post.PostID, post.GroupID, post.AuthorID, post.Content, post.CreatedAt)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}
