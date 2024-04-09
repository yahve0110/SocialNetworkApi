package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

type User struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Gender         string `json:"gender"`
	BirthDate      string `json:"birth_date"`
	ProfilePicture string `json:"profilePicture"`
	About          string `json:"about"`
}

// GetAllUsersExceptSubscribed returns all users except those the current user is subscribed to
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Get the user ID based on the current user's session
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get the user ID based on the current user's session
	currentUserID, err := helpers.GetUserIDFromSession(dbConnection, cookie.Value)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}


	// Execute the SQL query to get all users except the subscribed ones
	query := `
        SELECT user_id, username, first_name, last_name, gender, birth_date, profile_picture, about, email
        FROM users
        WHERE user_id NOT IN (
            SELECT user_followed
            FROM Followers
            WHERE user_following = ?
        )
    `
	rows, err := dbConnection.Query(query, currentUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing SQL query: %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate through the result set and build a slice of User structs
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Gender, &user.BirthDate, &user.ProfilePicture, &user.About, &user.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %s", err), http.StatusInternalServerError)
			return
		}

		users = append(users, user)
	}

	// Convert the slice of users to JSON
	usersJSON, err := json.Marshal(users)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(usersJSON)
}

// GetAllUsers is a handler to get all users
// func GetAllUsers(w http.ResponseWriter, r *http.Request) {
// 	// Access the global database connection from the db package
// 	dbConnection := database.DB

// 	// Execute the SQL query to get specific user fields
// 	rows, err := dbConnection.Query("SELECT user_id, username, first_name, last_name, gender, birth_date, profile_picture, about, email FROM users")
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Error executing SQL query: %s", err), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	// Iterate through the result set and build a slice of User structs
// 	var users []User
// 	for rows.Next() {
// 		var user User
// 		err := rows.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Gender, &user.BirthDate, &user.ProfilePicture, &user.About, &user.Email)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Error scanning row: %s", err), http.StatusInternalServerError)
// 			return
// 		}

// 		users = append(users, user)
// 	}

// 	// Convert the slice of users to JSON
// 	usersJSON, err := json.Marshal(users)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Set the Content-Type header and write the JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(usersJSON)
// }
