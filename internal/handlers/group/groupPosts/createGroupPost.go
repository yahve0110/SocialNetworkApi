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



	//upload group image to cloud storage
	postImageBase64 := postData.Image
	if postImageBase64 != "" {
		cloudinaryURL, err := helpers.ImageToCloud(postImageBase64, w)
		if err != nil {
			// Handle error
			return
		}
		postData.Image = cloudinaryURL
	}

	// Insert the post into the database
	err = InsertGroupPost(dbConnection, postData)
	if err != nil {
		log.Printf("Error inserting group post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Fetch the author's first name and last name from the database
	authorFirstName, authorLastName, err := GetUserFirstNameAndLastName(dbConnection, userID)
	if err != nil {
		log.Printf("Error fetching author's first name and last name: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Include the author's first name and last name in the response
	postData.AuthorFirstName = authorFirstName
	postData.AuthorLastName = authorLastName

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postData)
}

// InsertGroupPost inserts a new post into the database
func InsertGroupPost(db *sql.DB, post models.GroupPost) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO group_posts (post_id, group_id, author_id, content, post_date, group_post_img) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(post.PostID, post.GroupID, post.AuthorID, post.Content, post.CreatedAt,post.Image)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}

func GetUserFirstNameAndLastName(db *sql.DB, userID string) (string, string, error) {
    // Prepare the SQL statement to fetch the user's first name and last name
    query := "SELECT first_name, last_name FROM users WHERE user_id = ?"

    // Execute the SQL query to fetch the user's first name and last name
    row := db.QueryRow(query, userID)

    var firstName, lastName string
    // Scan the query result into variables
    err := row.Scan(&firstName, &lastName)
    if err != nil {
        return "", "", err
    }

    return firstName, lastName, nil
}