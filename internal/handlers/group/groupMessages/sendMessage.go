package groupChat

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
	"time"

	database "social/internal/db"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)


type GroupChatMessageToSave struct {
    MessageID  string    `json:"message_id"`
    Content    string    `json:"content"`
    AuthorID   string    `json:"author_id"`
    ChatID     string    `json:"chat_id"`
    CreatedAt  time.Time `json:"created_at"`
    AuthorName string    `json:"author_name"`
}

type SendMessageResponse struct {
	ChatID          string    `json:"chat_id"`
	MessageAuthorID string    `json:"author_id"`
	Content         string    `json:"content"`
	Timestamp       time.Time `json:"created_at"`
	FirstName       string    `json:"author_name"`
	LastName        string    `json:"last_name"`
	ProfilePicture  string    `json:"profile_picture"`
}

var (
	groupClients   = make(map[*clientInfo]bool)
	groupBroadcast = make(chan SendMessageResponse)
	groupClientsMu sync.Mutex
)

type MessageData struct {
	ChatID  string `json:"chat_id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type clientInfo struct {
	conn *websocket.Conn

	chatID string
}

type UserInfo struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
}

func HandleGroupChatConnections(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chatID")

	// Upgrade the HTTP connection to WebSocket
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ws.Close()

	dbConnection := database.DB

	// Register the new client with chatID
	groupClientsMu.Lock()
	groupClients[&clientInfo{conn: ws, chatID: chatID}] = true
	groupClientsMu.Unlock()

	for {
		var data MessageData
		// Read messages from WebSocket and decode into MessageData
		err := ws.ReadJSON(&data)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			groupClientsMu.Lock()
			for client := range groupClients {
				if client.conn == ws {
					delete(groupClients, client)
					break
				}
			}
			groupClientsMu.Unlock()
			break
		}

		msgTime := time.Now()
		err = saveGroupChatMessage(data, msgTime)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			break
		}

		authorInfo, err := getUserInfo(data.UserID, dbConnection)
		if err != nil {
			log.Printf("Error fetching user information: %v", err)
			break
		}

		// Формирование ответа
		response := SendMessageResponse{
			ChatID:          data.ChatID,
			Content:         data.Message,
			MessageAuthorID: data.UserID,
			Timestamp:       msgTime,
			FirstName:       authorInfo.FirstName,
			LastName:        authorInfo.LastName,
			ProfilePicture:  authorInfo.ProfilePicture,
		}

		groupClientsMu.Lock()
		for client := range groupClients {
			if client.chatID == data.ChatID {
				err := client.conn.WriteJSON(response)
				if err != nil {
					log.Printf("Ошибка записи: %v", err)
					client.conn.Close()
					delete(groupClients, client)
				}
			}
		}
		groupClientsMu.Unlock()
	}
}

func saveGroupChatMessage(message MessageData, timestamp time.Time) error {
	db := database.DB

	messageID := uuid.New().String()

	// Insert message into database
	_, err := db.Exec("INSERT INTO group_chat_messages (chat_id, author_id, content, created_at, message_id) VALUES (?, ?, ?, ?, ?)",
		message.ChatID, message.UserID, message.Message, timestamp, messageID)
	if err != nil {
		return err
	}

	return nil
}

func getUserInfo(userID string, db *sql.DB) (UserInfo, error) {
	var userInfo UserInfo
	err := db.QueryRow("SELECT first_name, last_name, profile_picture FROM users WHERE user_id = ?", userID).
		Scan(&userInfo.FirstName, &userInfo.LastName, &userInfo.ProfilePicture)
	if err != nil {
		return UserInfo{}, err
	}
	return userInfo, nil
}

func HandleGroupMessages() {

	for {
		msg := <-groupBroadcast

		groupClientsMu.Lock()
		for client := range groupClients {
			if client.chatID == msg.ChatID {
				err := client.conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error when writing: %v", err)
					client.conn.Close()
					delete(groupClients, client)
				}
			}
		}
		groupClientsMu.Unlock()
	}
}
