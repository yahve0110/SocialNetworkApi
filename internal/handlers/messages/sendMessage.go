package messageHandlers

import (
	"fmt"
	"log"
	"net/http"
	database "social/internal/db"
	"time"

	"github.com/gorilla/websocket"
)

type SendMessageRequest struct {
	ChatID  string `json:"chat_id"`
	Content string `json:"content"`
}

type SendMessageResponse struct {
	ChatID          string    `json:"chat_id"`
	MessageAuthorID string    `json:"message_author_id"`
	Content         string    `json:"content"`
	Timestamp       time.Time `json:"timestamp"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	ProfilePicture  string    `json:"profile_picture"`
}

// Структура для хранения сообщений
type MessageData struct {
	ChatID  string `json:"chat_id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// SendPrivateMessage обрабатывает запрос на отправку сообщения в чат
// func SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
// 	dbConnection := database.DB

// 	// Проверка аутентификации пользователя
// 	cookie, err := r.Cookie("sessionID")
// 	if err != nil {
// 		http.Error(w, "User not authenticated", http.StatusUnauthorized)
// 		return
// 	}
// 	userID, err := helpers.GetUserIDFromSession(dbConnection, cookie.Value)
// 	if err != nil {
// 		http.Error(w, "User not authenticated", http.StatusUnauthorized)
// 		return
// 	}

// 	// Получение данных о сообщении из тела запроса
// 	var requestData SendMessageRequest
// 	err = json.NewDecoder(r.Body).Decode(&requestData)
// 	if err != nil {
// 		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Вставка сообщения в базу данных
// 	result, err := dbConnection.Exec("INSERT INTO privatechat_messages (chat_id, message_author_id, content, timestamp) VALUES (?, ?, ?, ?)",
// 		requestData.ChatID, userID, requestData.Content, time.Now())
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Error inserting message into database: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Получение идентификатора нового сообщения
// 	messageID, err := result.LastInsertId()
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Error getting message ID: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Формирование ответа
// 	response := SendMessageResponse{
// 		MessageID: fmt.Sprintf("%d", messageID),
// 	}

// 	// Отправка ответа клиенту
// 	json.NewEncoder(w).Encode(response)
// }

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan SendMessageResponse)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to WebSocket
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ws.Close()

	dbConnection := database.DB

	// Register the new client
	clients[ws] = true

	for {
		var data MessageData
		// Read messages from WebSocket and decode into MessageData
		err := ws.ReadJSON(&data)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(clients, ws)
			break
		}

		msgTime := time.Now()
		// Вставка сообщения в базу данных
		// Insert the message into the database
		// Insert the message into the database
		_, err = dbConnection.Exec("INSERT INTO privatechat_messages (chat_id, message_author_id, content, timestamp) VALUES (?, ?, ?, ?)",
			data.ChatID, data.UserID, data.Message, msgTime)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting message into database: %v", err), http.StatusInternalServerError)
			return
		}

		var firstName, lastName, profilePicture string
		err = dbConnection.QueryRow("SELECT first_name, last_name, profile_picture FROM users WHERE user_id = ?", data.UserID).Scan(&firstName, &lastName, &profilePicture)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching user information: %v", err), http.StatusInternalServerError)
			return
		}

		// Формирование ответа
		response := SendMessageResponse{
			ChatID:          data.ChatID,
			Content:         data.Message,
			MessageAuthorID: data.UserID,
			Timestamp:       msgTime,
			FirstName:       firstName,
			LastName:        lastName,
			ProfilePicture:  profilePicture,
		}

		// Print the received message to the console
	
		broadcast <- response
	}
}

func HandleMessages() {
	for {
		// Получение сообщения из общего канала
		msg := <-broadcast
		// Отправка сообщения всем клиентам
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Ошибка записи: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
