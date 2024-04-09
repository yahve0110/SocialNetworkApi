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

// GetAllGroupPostsHandler handles the retrieval of all group posts

func GetAllGroupPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Check if the user is authenticated (you may modify this based on your authentication mechanism)
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

	// Get the group ID from the request (you need to define how the group ID is sent in the request)
	groupID := r.FormValue("group_id")

	// Check if the user is a member of the group (you may modify this based on your membership verification)
	isGroupMember, err := helpers.IsUserGroupMember(dbConnection, userID, groupID)
	if err != nil {
		log.Printf("Error checking if user is a group member: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	isGroupCreator, err := helpers.IsUserGroupCreator(dbConnection, userID, groupID)
	if err != nil {
		log.Printf("Error checking if user is a group creator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("group cretor: ", isGroupCreator)

	if !isGroupMember && !isGroupCreator {
		http.Error(w, "Unauthorized: User is not a member of the group", http.StatusUnauthorized)
		return
	}

	// Retrieve all group posts for the group
	groupPosts, err := getAllGroupPostsFromDatabase(dbConnection, groupID)
	if err != nil {
		log.Printf("Error fetching group posts from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Iterate through the group posts and fetch the count of likes for each post
	for i := range groupPosts {
		likeCount, err := getLikeCountForPost(dbConnection, groupPosts[i].PostID)
		if err != nil {
			log.Printf("Error fetching like count for post %s: %v", groupPosts[i].PostID, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		groupPosts[i].LikeCount = likeCount
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupPosts)
}

// getAllGroupPostsFromDatabase fetches all group posts for a group from the database
func getAllGroupPostsFromDatabase(dbConnection *sql.DB, groupID string) ([]models.GroupPost, error) {
	// Query all group posts for the group from the "group_posts" table
	rows, err := dbConnection.Query("SELECT post_id, group_id, author_id, content, post_date FROM group_posts WHERE group_id = ?", groupID)
	if err != nil {
		log.Printf("Error querying group posts from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and create group post objects
	var groupPosts []models.GroupPost
	for rows.Next() {
		var groupPost models.GroupPost
		if err := rows.Scan(&groupPost.PostID, &groupPost.GroupID, &groupPost.AuthorID, &groupPost.Content, &groupPost.CreatedAt); err != nil {
			log.Printf("Error scanning group post rows: %v", err)
			return nil, err
		}
		groupPosts = append(groupPosts, groupPost)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over group post rows: %v", err)
		return nil, err
	}

	return groupPosts, nil
}

// getLikeCountForPost fetches the count of likes for a group post
func getLikeCountForPost(dbConnection *sql.DB, postID string) (int, error) {
	// Query the group_post_likes table to get the count of likes for the post
	query := "SELECT COUNT(*) FROM group_post_likes WHERE post_id = ?"
	var likeCount int
	err := dbConnection.QueryRow(query, postID).Scan(&likeCount)
	if err != nil {
		log.Printf("Error fetching like count for post %s: %v", postID, err)
		return 0, err
	}

	return likeCount, nil
}
