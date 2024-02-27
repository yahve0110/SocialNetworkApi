package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
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

// GetAllUsers is a handler to get all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Execute the SQL query to get specific user fields
	rows, err := dbConnection.Query("SELECT user_id, username, first_name, last_name, gender, birth_date, profile_picture, about, email FROM users")
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
