package commentHandlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

// Modify the DeleteComment function
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	// Extract the comment ID from the request
	commentID := r.FormValue("comment_id")

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Check if the user is authenticated
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

	// Check if the user is the creator of the comment
	isCommentCreator, err := IsUserCommentCreator(dbConnection, userID, commentID)
	if err != nil {
		fmt.Println("Error checking if user is the comment creator:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isCommentCreator {
		http.Error(w, "Unauthorized: User is not the creator of the comment", http.StatusUnauthorized)
		return
	}

	imageURL, err := GetCommentImageURLFromDatabase(dbConnection, commentID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if imageURL != "" {
		err = helpers.DeleteFromCloudinary(imageURL)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// If the user is the creator, proceed to delete the comment
	_, err = dbConnection.Exec("DELETE FROM comments WHERE comment_id = ?", commentID)
	if err != nil {
		fmt.Println("Error deleting comment from database:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Comment deleted successfully"))
}

// IsUserCommentCreator checks if the user with the given userID is the creator of the comment with the given commentID
func IsUserCommentCreator(db *sql.DB, userID, commentID string) (bool, error) {
	// Query the comments table to check if the user is the creator of the comment
	query := "SELECT EXISTS(SELECT 1 FROM comments WHERE author_id = ? AND comment_id = ?)"
	var exists bool
	err := db.QueryRow(query, userID, commentID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if user is the comment creator: %v", err)
		return false, err
	}

	return exists, nil
}

func GetCommentImageURLFromDatabase(db *sql.DB, postID string) (string, error) {
	// Подготовьте SQL-запрос для извлечения URL изображения из базы данных
	query := "SELECT image FROM comments WHERE comment_id  = ?"

	// Выполните запрос к базе данных
	var imageURL string
	err := db.QueryRow(query, postID).Scan(&imageURL)
	if err != nil {
		return "", err
	}

	return imageURL, nil

}
