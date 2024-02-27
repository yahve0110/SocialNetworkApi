package followHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

type FollowRequest struct {
	UserFollowing  string
	UserFollowed   string `json:"user_followed"`
	UserUnfollowed string `json:"user_to_unfollow"`
	UserPending string `json:"user_pending"`
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Parse the request body to get the follow request
	var followRequest FollowRequest
	err := json.NewDecoder(r.Body).Decode(&followRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Retrieve session ID from cookies
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		// Cookie not found, session is invalid
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID := cookie.Value

	userFollowing, err := helpers.GetUserIDFromSession(dbConnection, sessionID)
	if err != nil {
		http.Error(w, "Error getting user ID from session", http.StatusInternalServerError)
		return
	}
	followRequest.UserFollowing = userFollowing

	// Check if the users exist
	userFollowingExists, err := helpers.UserExists(dbConnection, followRequest.UserFollowing)
	if err != nil {
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}
	userFollowedExists, err := helpers.UserExists(dbConnection, followRequest.UserFollowed)
	if err != nil {
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}

	if !userFollowingExists || !userFollowedExists {
		http.Error(w, "User following or user followed not found", http.StatusNotFound)
		return
	}

	// Check if the user being followed has a public profile
	privacyStatus, err := helpers.IsUserProfilePublic(dbConnection, followRequest.UserFollowed)
	if err != nil {
		http.Error(w, "Error checking user profile privacy", http.StatusInternalServerError)
		return
	}

	// Insert or update the follow status based on the profile privacy
	if privacyStatus == "public" {
		// The profile is public, automatically add the follower
		err = helpers.InsertOrUpdateFollowStatus(dbConnection, followRequest.UserFollowing, followRequest.UserFollowed)
		if err != nil {
			http.Error(w, "Error updating follow status", http.StatusInternalServerError)
			return
		}

		// Respond with a success message for automatic follow
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User followed successfully"))
	} else {
		// The profile is private, add a follow request
		err = helpers.InsertFollowRequest(dbConnection, followRequest.UserFollowing, followRequest.UserFollowed)
		if err != nil {
			http.Error(w, "Error adding follow request", http.StatusInternalServerError)
			return
		}

		// Respond with a success message for follow request
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Follow request sent successfully"))
	}
}

// UnfollowUserHandler is a handler to unfollow a user
func UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Retrieve session ID from cookies
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		// Cookie not found, session is invalid
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID := cookie.Value

	// Get the user ID based on the current user's session
	userFollowing, err := helpers.GetUserIDFromSession(dbConnection, sessionID)
	if err != nil {
		http.Error(w, "Error getting user ID from session", http.StatusInternalServerError)
		return
	}

	// Parse the unfollow request from the request body
	var unfollowRequest FollowRequest
	err = json.NewDecoder(r.Body).Decode(&unfollowRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the users exist
	userFollowingExists, err := helpers.UserExists(dbConnection, userFollowing)
	if err != nil {
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}
	userUnfollowedExists, err := helpers.UserExists(dbConnection, unfollowRequest.UserUnfollowed)
	if err != nil {
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}

	if !userFollowingExists || !userUnfollowedExists {
		http.Error(w, "User following or user to unfollow not found", http.StatusNotFound)
		return
	}

	// Delete the follow status in the Followers table
	err = helpers.DeleteFollowStatus(dbConnection, userFollowing, unfollowRequest.UserUnfollowed)
	if err != nil {
		http.Error(w, "Error deleting follow status", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Unfollowed successfully"))
}

// AcceptPendingFollowerHandler is a handler to accept a pending follower request
func AcceptPendingFollowerHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Retrieve session ID from cookies
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		// Cookie not found, session is invalid
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID := cookie.Value

	// Get the user ID based on the current user's session
	userFollowed, err := helpers.GetUserIDFromSession(dbConnection, sessionID)
	if err != nil {
		http.Error(w, "Error getting user ID from session", http.StatusInternalServerError)
		return
	}

	fmt.Println("userFollowed:",userFollowed)

	// Parse the accept request from the request body
	var acceptRequest FollowRequest
	err = json.NewDecoder(r.Body).Decode(&acceptRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("userFollowing:",acceptRequest.UserPending)
	// Check if the users exist


	// Accept the pending follower request
	err = helpers.AcceptPendingFollower(dbConnection, acceptRequest.UserPending,userFollowed)
	if err != nil {
		http.Error(w, "Error accepting pending follower request", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pending follower request accepted successfully"))
}
