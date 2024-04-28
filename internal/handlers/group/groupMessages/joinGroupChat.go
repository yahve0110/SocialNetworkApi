package groupChat

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

type JoinGroupChatRequest struct {
	GroupID string `json:"group_id"`
}

func JoinGroupChatHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a JoinGroupChatRequest struct
	var joinRequest JoinGroupChatRequest
	if err := json.NewDecoder(r.Body).Decode(&joinRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the user ID based on the current user's session
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Access the global database connection
	dbConnection := database.DB

	// Get the user ID based on the current user's session
	userID, err := helpers.GetUserIDFromSession(dbConnection, cookie.Value)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Validate that required fields are not empty
	if joinRequest.GroupID == "" {
		http.Error(w, "GroupID is required", http.StatusBadRequest)
		return
	}

	// Insert the user into the group chat members
	if err := addUserToGroupChat(userID, joinRequest.GroupID, dbConnection); err != nil {
		log.Printf("Error adding user to group chat: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User joined group chat successfully"})
}

// addUserToGroupChat добавляет нового пользователя в таблицу участников группового чата
func addUserToGroupChat(userID, groupID string, db *sql.DB) error {
	// Проверяем, существует ли уже запись для данного пользователя и группы
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM group_chat_members WHERE member_id = ? AND chat_id = ?", userID, groupID).Scan(&count)
	if err != nil {
		return err
	}

	// Если запись уже существует, возвращаем ошибку
	if count > 0 {
		return errors.New("user already joined this group chat")
	}

	// Подготовка SQL-запроса
	stmt, err := db.Prepare("INSERT INTO group_chat_members (member_id, chat_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение SQL-запроса
	_, err = stmt.Exec(userID, groupID)
	if err != nil {
		return err
	}

	return nil
}
