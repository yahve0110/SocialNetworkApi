package groupChat

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	database "social/internal/db"
	"social/internal/helpers"

	"github.com/google/uuid"
)

type GroupChatMessageToSave struct {
	Content   string `json:"content"`
	AuthorID  string `json:"author_id"`
	ChatID    string `json:"chat_id"`
	CreatedAt string `json:"created_at"`
}

func SendGroupChatMessage(w http.ResponseWriter, r *http.Request) {
	// Декодирование данных сообщения из тела запроса
	var message GroupChatMessageToSave
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dbConnection := database.DB
	// Получение идентификатора текущего пользователя из сессии
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	userID, err := helpers.GetUserIDFromSession(dbConnection, cookie.Value)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	// Получение дополнительных данных из запроса, таких как chat_id
	chatID := message.ChatID // Используем chatID из тела запроса
	authorID := userID

	// Установка времени создания сообщения
	createdAt := time.Now().Format(time.RFC3339)
	messageId := uuid.New().String()
	// Сохранение сообщения в базе данных
	err = saveGroupChatMessage(dbConnection, message.Content, chatID, authorID, createdAt, messageId)
	if err != nil {
		http.Error(w, "Failed to send group chat message", http.StatusInternalServerError)
		return
	}

	// Отправка ответа клиенту об успешной отправке сообщения
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Group chat message sent successfully"})
}

func saveGroupChatMessage(db *sql.DB, message, chatID, authorID, createdAt, message_id string) error {
	// Подготовка SQL-запроса для вставки сообщения в базу данных
	stmt, err := db.Prepare("INSERT INTO group_chat_messages (content, author_id, chat_id, created_at,message_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение SQL-запроса
	_, err = stmt.Exec(message, authorID, chatID, createdAt,message_id)
	if err != nil {
		return err
	}

	return nil
}
