package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"

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
	err = helpers.ValidateCredentials(dbConnection, newUser.Email, newUser.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser.Role = "user"

	// Insert user into the database
	err = helpers.InsertUser(dbConnection, newUser)
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
