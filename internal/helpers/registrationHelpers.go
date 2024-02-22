package helpers

import (
	"database/sql"
	"fmt"
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
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE nickname = ?", nickname).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}