package userHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// Get the database connection
	dbConnection := database.DB

	// Get the user ID based on the current user's session
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

	// Execute the SQL query to get specific user fields
	rows, err := dbConnection.Query("SELECT user_id, username, first_name, last_name, gender, birth_date, profile_picture, about, email FROM users WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing SQL query: %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Process the query result
	var user models.User
	for rows.Next() {
		err := rows.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Gender, &user.BirthDate, &user.ProfilePicture, &user.About, &user.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %s", err), http.StatusInternalServerError)
			return
		}

		// Send user data as JSON response
		jsonData, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}

	// If no rows are returned, it means user not found
	http.Error(w, "User not found", http.StatusNotFound)
}
