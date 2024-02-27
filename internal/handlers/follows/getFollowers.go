package followHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

// GetFollowersHandler is a handler to get all followers for a given user
func GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Parse the user ID from the request parameters or headers (adjust accordingly)
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get followers for the specified user
	followers, err := helpers.GetFollowers(dbConnection, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting followers: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert the slice of followers to JSON
	followersJSON, err := json.Marshal(followers)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(followersJSON)
}

// GetFollowingHandler is a handler to get all users a given user is following
func GetFollowingHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Parse the user ID from the request parameters or headers (adjust accordingly)
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get followed users for the specified user
	following, err := helpers.GetFollowing(dbConnection, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting followed users: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert the slice of followed users to JSON
	followingJSON, err := json.Marshal(following)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(followingJSON)
}

// GetFollowersWithPendingStatusHandler is a handler to get followers with a 'pending' status
func GetFollowersWithPendingStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Parse the user ID from the request parameters or headers (adjust accordingly)
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get followers with a 'pending' status for the specified user
	followers, err := helpers.GetFollowersWithStatus(dbConnection, userID, "pending")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting followers with pending status: %v", err), http.StatusInternalServerError)
		return
	}

	// Get detailed information about each follower
	var detailedFollowers []models.User
	for _, followerID := range followers {
		follower, err := helpers.GetUserByID(dbConnection, followerID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting follower information: %v", err), http.StatusInternalServerError)
			return
		}
		detailedFollowers = append(detailedFollowers, follower)
	}

	// Convert the slice of detailed followers to JSON
	followersJSON, err := json.Marshal(detailedFollowers)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(followersJSON)
}
