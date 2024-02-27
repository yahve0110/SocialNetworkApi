package postHandler

import (
	"encoding/json"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

type PostLike struct {
	PostID string `json:"post_id"`
}

func AddPostLike(w http.ResponseWriter, r *http.Request) {
    // Access the global database connection from the db package
    dbConnection := database.DB

    // Get user ID from session (assuming session ID is stored in a cookie)
    sessionID, err := r.Cookie("sessionID")
    if err != nil {
        http.Error(w, "Failed to get session ID from cookie", http.StatusBadRequest)
        return
    }

    userID, err := helpers.GetUserIDFromSession(dbConnection, sessionID.Value)
    if err != nil {
        http.Error(w, "Failed to get user ID from session: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Parse the request body to get post-like data
    var postLike PostLike
    err = json.NewDecoder(r.Body).Decode(&postLike)
    if err != nil {
        http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Check if the user has already liked the post
    var existingUserID string
    err = dbConnection.QueryRow("SELECT user_id FROM postLikes WHERE post_id = ? AND user_id = ?", postLike.PostID, userID).Scan(&existingUserID)
    if err == nil {
        // If the user has already liked the post, remove the like
        _, err = dbConnection.Exec("DELETE FROM postLikes WHERE post_id = ? AND user_id = ?", postLike.PostID, userID)
    } else {
        // If the user hasn't liked the post, add the like
        _, err = dbConnection.Exec("INSERT INTO postLikes (post_id, user_id) VALUES (?, ?)", postLike.PostID, userID)
    }

    if err != nil {
        http.Error(w, "Failed to add/remove post like: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.WriteHeader(http.StatusOK)
}
