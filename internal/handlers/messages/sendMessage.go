package messageHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"time"
)

type SendMessageRequest struct {
	ChatID  string `json:"chat_id"`
	Content string `json:"content"`
}

type SendMessageResponse struct {
	MessageID string `json:"message_id"`
}

// SendPrivateMessage обрабатывает запрос на отправку сообщения в чат
func SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	dbConnection := database.DB

	// Проверка аутентификации пользователя
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

	// Получение данных о сообщении из тела запроса
	var requestData SendMessageRequest
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Вставка сообщения в базу данных
	result, err := dbConnection.Exec("INSERT INTO privatechat_messages (chat_id, message_author_id, content, timestamp) VALUES (?, ?, ?, ?)",
		requestData.ChatID, userID, requestData.Content, time.Now())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting message into database: %v", err), http.StatusInternalServerError)
		return
	}

	// Получение идентификатора нового сообщения
	messageID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting message ID: %v", err), http.StatusInternalServerError)
		return
	}

	// Формирование ответа
	response := SendMessageResponse{
		MessageID: fmt.Sprintf("%d", messageID),
	}

	// Отправка ответа клиенту
	json.NewEncoder(w).Encode(response)
}
