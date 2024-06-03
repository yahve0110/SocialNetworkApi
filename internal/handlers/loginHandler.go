package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
	"strings"
	"time"
)

type CookieData struct {
	Name     string
	Value    string
	Expires  time.Time
	HttpOnly bool
	SameSite http.SameSite
	Secure   bool
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	// Decode the JSON request body into the Credentials struct
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dbConnection := database.DB

	if strings.Contains(credentials.Username, "@") {
		// Retrieve the username from the database using the provided email
		username, err := helpers.GetUsernameByEmail(dbConnection, credentials.Username)
		if err != nil {
			http.Error(w, "Invalid email", http.StatusUnauthorized)
			return
		}
		credentials.Username = username
	}

	// Check if the username exists in the database
	userExists, err := helpers.IsUsernameExists(dbConnection, credentials.Username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error checking username"})
		return
	}

	if !userExists {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}

	// Retrieve hashed password from the database
	storedPassword, err := helpers.GetPasswordByUsername(dbConnection, credentials.Username)
	if err != nil {
		http.Error(w, "Error retrieving password", http.StatusInternalServerError)
		return
	}

	// Compare the provided password with the stored hash
	if err := helpers.CheckPasswordHash(credentials.Password, storedPassword); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	cookie, err := helpers.CreateSession(w, dbConnection, credentials.Username)
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	// Get the user ID
	userID, err := helpers.GetUserIDByUsername(dbConnection, credentials.Username)
	if err != nil {
		http.Error(w, "Error retrieving user ID", http.StatusInternalServerError)
		return
	}

	// Successful login actions
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"cookie":  cookie,
		"userID":  userID,
	})
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	dbConnection := database.DB
	err := helpers.Logout(w, r, dbConnection)
	if err != nil {
		http.Error(w, "Error during logout", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}
