package postHandler

import (
	"encoding/json"

	"net/http"
	"social/internal/db"
	"social/internal/helpers"
	"database/sql"
)

type PostLike struct {
	PostID string `json:"post_id"`
}

func AddPostLike(w http.ResponseWriter, r *http.Request) {
	dbConnection := database.DB

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

	var postLike PostLike
	err = json.NewDecoder(r.Body).Decode(&postLike)
	if err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var existingUserID string
	err = dbConnection.QueryRow("SELECT user_id FROM postLikes WHERE post_id = ? AND user_id = ?", postLike.PostID, userID).Scan(&existingUserID)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Failed to check existing post like: "+err.Error(), http.StatusInternalServerError)
		return
	}


    if err != sql.ErrNoRows {
        _, err = dbConnection.Exec("DELETE FROM postLikes WHERE post_id = ? AND user_id = ?", postLike.PostID, userID)
        if err != nil {
            http.Error(w, "Failed to remove post like: "+err.Error(), http.StatusInternalServerError)
            return
        }
    } else {
        _, err = dbConnection.Exec("INSERT INTO postLikes (post_id, user_id) VALUES (?, ?)", postLike.PostID, userID)
        if err != nil {
            http.Error(w, "Failed to add post like: "+err.Error(), http.StatusInternalServerError)
            return
        }
    }


	var likesCount int
	err = dbConnection.QueryRow("SELECT COUNT(*) FROM postLikes WHERE post_id = ?", postLike.PostID).Scan(&likesCount)
	if err != nil {
		http.Error(w, "Failed to get post likes count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	responseData := map[string]interface{}{
		"likes_count": likesCount,
	}
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Failed to encode response data: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
