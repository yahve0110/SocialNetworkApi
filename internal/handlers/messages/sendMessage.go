package messageHandlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	database "social/internal/db"

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

type MessageData struct {
	ChatID  string `json:"chat_id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type clientInfo struct {
	conn   *websocket.Conn
	userID string
	chatID string
}

var (
	clients   = make(map[*clientInfo]bool)
	broadcast = make(chan SendMessageResponse)
	clientsMu sync.Mutex
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	chatID := r.URL.Query().Get("chatID")

	// Upgrade the HTTP connection to WebSocket
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ws.Close()

	dbConnection := database.DB
	fmt.Println("USER ID:", userID)

	// Register the new client with chatID
	clientsMu.Lock()
	clients[&clientInfo{conn: ws, userID: userID, chatID: chatID}] = true
	clientsMu.Unlock()

	for {
		var data MessageData
		// Read messages from WebSocket and decode into MessageData
		err := ws.ReadJSON(&data)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			clientsMu.Lock()
			for client := range clients {
				if client.conn == ws {
					delete(clients, client)
					break
				}
			}
			clientsMu.Unlock()
			break
		}

		msgTime := time.Now()
		_, err = dbConnection.Exec("INSERT INTO privatechat_messages (chat_id, message_author_id, content, timestamp) VALUES (?, ?, ?, ?)",
			data.ChatID, data.UserID, data.Message, msgTime)
		if err != nil {
			log.Printf("Error inserting message into database: %v", err)
			break
		}

		var firstName, lastName, profilePicture string
		err = dbConnection.QueryRow("SELECT first_name, last_name, profile_picture FROM users WHERE user_id = ?", data.UserID).Scan(&firstName, &lastName, &profilePicture)
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
			FirstName:       firstName,
			LastName:        lastName,
			ProfilePicture:  profilePicture,
		}

		clientsMu.Lock()
		for client := range clients {
			if client.chatID == data.ChatID {
				err := client.conn.WriteJSON(response)
				if err != nil {
					log.Printf("Websocket writing error: %v", err)
					client.conn.Close()
					delete(clients, client)
				}
			}
		}
		clientsMu.Unlock()
	}
}

func HandleMessages() {
	dbConnection := database.DB

	for {
		msg := <-broadcast

		participants, err := GetChatParticipants(msg.ChatID, dbConnection)
		if err != nil {
			log.Printf("Error getting chat participants: %v", err)
			continue
		}

		clientsMu.Lock()
		for client := range clients {
			for _, participant := range participants {
				if participant == client.userID && participant != msg.MessageAuthorID {
					err := client.conn.WriteJSON(msg)
					if err != nil {
						log.Printf("Ошибка записи: %v", err)
						client.conn.Close()
						delete(clients, client)
					}
				}
			}
		}
		clientsMu.Unlock()
	}
}

func GetChatParticipants(chatID string, db *sql.DB) ([]string, error) {
	participants := make([]string, 0)

	rows, err := db.Query("SELECT user1_id, user2_id FROM privatechat WHERE chat_id = ?", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user1ID, user2ID string
		if err := rows.Scan(&user1ID, &user2ID); err != nil {
			return nil, err
		}
		participants = append(participants, user1ID, user2ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return participants, nil
}
