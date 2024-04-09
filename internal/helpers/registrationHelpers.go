package helpers

import (
	"database/sql"
	"fmt"
	"social/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// ValidateCredentials checks if the provided email and nickname are unique
func ValidateCredentials(db *sql.DB, email, nickname string) error {
	emailUnique, err := IsEmailUnique(db, email)
	if err != nil {
		return err
	}
	if !emailUnique {
		return fmt.Errorf("invalid credentials")
	}

	nicknameUnique, err := IsNicknameUnique(db, nickname)
	if err != nil {
		return err
	}
	if !nicknameUnique {
		return fmt.Errorf("invalid credentials")
	}

	return nil
}

// IsEmailUnique checks if the given email is unique in the database
func IsEmailUnique(db *sql.DB, email string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// IsNicknameUnique checks if the given nickname is unique in the database
func IsNicknameUnique(db *sql.DB, nickname string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", nickname).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}


 //insert user in database
func InsertUser(dbConnection *sql.DB, user models.User) error {
	fmt.Println(user)
	if user.Email == "" || user.Password == ""  || user.Username == ""{
		return fmt.Errorf("email and password are required")
	}
	_, err := dbConnection.Exec(`
		INSERT INTO users (user_id, Username, first_name, last_name, email, password, gender, birth_date, profile_picture, role,about,privacy)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, user.UserID, user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.Gender, user.BirthDate, user.ProfilePicture, user.Role,user.About,user.Privacy)
	return err
}
