package userHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/models"
)

func GetUserInfoById(w http.ResponseWriter, r *http.Request) {
	// Extract user_id from request URL parameters
	userID := r.URL.Query().Get("user_id")

	// Get the database connection
	dbConnection := database.DB

	// Execute the SQL query to get specific user fields
	rows, err := dbConnection.Query("SELECT user_id, username, first_name, last_name, gender, birth_date, profile_picture, about, email, privacy FROM users WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing SQL query: %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Check if any rows were found
	if !rows.Next() {
		http.Error(w, "No rows found for user ID", http.StatusNotFound)
		return
	}

	// Process the query result
	var user models.User
	err = rows.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Gender, &user.BirthDate, &user.ProfilePicture, &user.About, &user.Email,  &user.Privacy)
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
}
