package helpers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)
type CookieData struct {
	Name     string
	Value    string
	Expires  time.Time
	HttpOnly bool
	SameSite http.SameSite
	Secure   bool
}

// IsUsernameExists checks if the given username exists in the database
func IsUsernameExists(db *sql.DB, username string) (bool, error) {
	var count int
	fmt.Println("username: ", username)
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking username: %v", err)
	}
	return count > 0, nil
}

// GetPasswordByUsername retrieves the hashed password for the given username from the database
func GetPasswordByUsername(db *sql.DB, username string) (string, error) {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("username not found")
		}
		return "", fmt.Errorf("error retrieving password: %v", err)
	}
	return hashedPassword, nil
}

// CheckPasswordHash compares a plain-text password with a hashed password
func CheckPasswordHash(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err // Passwords don't match
	}
	return nil // Passwords match
}

//sessions

// SaveSessionInfo saves the session information in the database
func SaveSessionInfo(db *sql.DB, userID, sessionID string, expirationTime time.Time) error {
	fmt.Println("tuta")
	fmt.Println(sessionID, userID)

	expirationTimeString := expirationTime.Format("2006-01-02 15:04:05")

	// Assuming you have a "sessions" table with columns "user_id" and "session_id"
	_, err := db.Exec(`
   INSERT INTO sessions (user_id, session_id, expiration_time)
   VALUES (?, ?, ?)
`, userID, sessionID, expirationTimeString)

	if err != nil {
		fmt.Println("Error inserting session information into the database:", err)
		return fmt.Errorf("error saving session information: %v", err)
	}

	return nil
}

func CreateSession(w http.ResponseWriter, db *sql.DB, username string) (*http.Cookie, error) {
	// Generate a UUID for the session ID
	sessionID := uuid.New().String()

	// Get the user ID based on the username
	var userID string
	err := db.QueryRow("SELECT user_id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user ID for username %s: %v", username, err)
	}

	// Set a cookie with the session ID
	expiration := time.Now().Add(24 * time.Hour) // Adjust the expiration time as needed
	cookie := &http.Cookie{
		Name:     "sessionID",
		Value:    sessionID,
		Expires:  expiration,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	// Save session information in the database
	err = SaveSessionInfo(db, userID, sessionID, expiration)
	if err != nil {
		return nil, err
	}

	// Set the cookie in the response
	http.SetCookie(w, cookie)

	// Возвращаем данные куки
	return cookie, nil
}

func Logout(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	// Retrieve the sessionID from the cookie
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		// No session cookie found, consider the user as already logged out
		return nil
	}

	// Delete the session information from the database
	sessionID := cookie.Value
	_, err = db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	if err != nil {
		return fmt.Errorf("error deleting session information: %v", err)
	}

	// Expire the session cookie
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)

	return nil
}
