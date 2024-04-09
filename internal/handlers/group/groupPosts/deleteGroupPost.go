package groupPostHandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

// DeleteGroupPostHandler handles the deletion of a group post
func DeleteGroupPostHandler(w http.ResponseWriter, r *http.Request) {
	var requestData models.GroupPost

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Check if the requester is the group creator or the post creator
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

	// Check if the user is the group creator or the post creator
	isGroupCreator, err := helpers.IsUserGroupCreator(dbConnection, userID, requestData.GroupID)
	if err != nil {
		log.Printf("Error checking group creator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	isPostCreator, err := helpers.IsUserPostCreator(dbConnection, userID, requestData.PostID)
	if err != nil {
		log.Printf("Error checking post creator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isGroupCreator && !isPostCreator {
		http.Error(w, "Unauthorized: Only group creator or post creator can delete the post", http.StatusUnauthorized)
		return
	}
	fmt.Println("postID: ", requestData.PostID)

	// Delete the group post
	err = DeleteGroupPost(dbConnection, userID, requestData.PostID)
	if err != nil {
		log.Printf("Error deleting group post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestData)
}

// DeleteGroupPost deletes a post from the database
func DeleteGroupPost(db *sql.DB, userID, postID string) error {
	// Check if the user is the creator of the post
	// Prepare the SQL statement
	stmt, err := db.Prepare("DELETE FROM group_posts WHERE post_id = ? AND author_id = ?")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(postID, userID)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}
