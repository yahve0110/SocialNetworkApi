package postHandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

// DeletePostHandler handles the deletion of a post
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	var requestData models.Post

	// Decode the request body into a Post struct
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the database package
	dbConnection := database.DB

	// Check if the requester is authenticated
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

	// Check if the user is the creator of the post
	isPostCreator, err := IsUserPostCreator(dbConnection, userID, requestData.PostID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isPostCreator {
		http.Error(w, "Unauthorized: Only the post creator can delete the post", http.StatusUnauthorized)
		return
	}

	imageURL, err := GetImageURLFromDatabase(dbConnection, requestData.PostID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("img url to delete: ", imageURL)

	// Если imageURL не пустой, удалите изображение из облачного хранилища
	if imageURL != "" {
		err = helpers.DeleteFromCloudinary(imageURL)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Удалить все комментарии к этому посту
	err = DeleteCommentsByPostID(dbConnection, requestData.PostID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Call the function to delete the post from the database
	err = DeletePost(dbConnection, requestData.PostID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post deleted successfully"))
}

// IsUserPostCreator checks if a user is the creator of the specified post
func IsUserPostCreator(db *sql.DB, userID, postID string) (bool, error) {
	// Query the posts table to check if the user is the creator
	query := "SELECT EXISTS(SELECT 1 FROM posts WHERE author_id = ? AND post_id = ?)"
	var exists bool
	err := db.QueryRow(query, userID, postID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// DeletePost deletes a post from the database
func DeletePost(db *sql.DB, postID string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("DELETE FROM posts WHERE post_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(postID)
	if err != nil {
		return err
	}

	return nil
}

func GetImageURLFromDatabase(db *sql.DB, postID string) (string, error) {
	// Подготовьте SQL-запрос для извлечения URL изображения из базы данных
	query := "SELECT image FROM posts WHERE post_id = ?"

	// Выполните запрос к базе данных
	var imageURL string
	err := db.QueryRow(query, postID).Scan(&imageURL)
	if err != nil {
		return "", err
	}

	// Верните извлеченный URL изображения
	return imageURL, nil

}

// DeleteCommentsByPostID deletes all comments related to a post from the database
func DeleteCommentsByPostID(db *sql.DB, postID string) error {
    // Prepare SQL statement to select image URLs of comments related to the specified post
    stmt, err := db.Prepare("SELECT image FROM comments WHERE post_id = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Execute the SQL statement to retrieve image URLs
    rows, err := stmt.Query(postID)
    if err != nil {
        return err
    }
    defer rows.Close()

    // Iterate over the rows to retrieve image URLs and delete them from Cloudinary
    for rows.Next() {
        var imageURL string
        if err := rows.Scan(&imageURL); err != nil {
            return err
        }

        // If imageURL is not empty, delete the image from Cloudinary
        if imageURL != "" {
            if err := helpers.DeleteFromCloudinary(imageURL); err != nil {
				return err
            }
        }
    }

    // Check for any errors encountered while iterating over rows
    if err := rows.Err(); err != nil {
        return err
    }

    // Prepare SQL statement to delete comments related to the specified post
    deleteStmt, err := db.Prepare("DELETE FROM comments WHERE post_id = ?")
    if err != nil {
        return err
    }
    defer deleteStmt.Close()

    // Execute the SQL statement to delete comments
    _, err = deleteStmt.Exec(postID)
    if err != nil {
        return err
    }

    return nil
}
