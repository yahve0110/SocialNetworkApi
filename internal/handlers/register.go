package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social/internal/helpers"
	"social/internal/models"
	"social/internal/db"
	"github.com/google/uuid"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	fmt.Println("entered")

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate unique UserId
	newUser.UserID = uuid.New().String()

	// Hash password
	hashedPassword, err := helpers.HashPassword(newUser.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Assign hashed password to user
	newUser.Password = string(hashedPassword)

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Validate email and nickname uniqueness
	err = helpers.ValidateCredentials(dbConnection, newUser.Email, newUser.Nickname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
fmt.Println(newUser)
	// Insert user into the database
	_, err = dbConnection.Exec(`
		INSERT INTO users (user_id, nickname, first_name, last_name, email, password, gender, birth_date, profile_picture, role)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, newUser.UserID, newUser.Nickname, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Password, newUser.Gender, newUser.BirthDate, newUser.ProfilePicture, newUser.Role)
	if err != nil {
		// Log the error for debugging purposes
		fmt.Println("Error inserting user into database:", err)
		// Return an HTTP response with a 500 Internal Server Error status
		http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
