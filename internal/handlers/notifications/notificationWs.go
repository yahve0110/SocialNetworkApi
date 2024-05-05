package notifications

import (
	"fmt"
	"log"
	"net/http"
	database "social/internal/db"
	groupHandlers "social/internal/handlers/group"
	"sync"

	"github.com/gorilla/websocket"
)

type clientInfo struct {
	conn   *websocket.Conn
	userID string
}

type SendMessageResponse struct {
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

var (
	clients   = make(map[*clientInfo]bool)
	broadcast = make(chan SendMessageResponse)
	clientsMu sync.Mutex // мьютекс для синхронизации доступа к clients
)

type MessageData struct {
	UserID   string `json:"user_id"`
	SenderID string `json:"sender_id"`
	Message  string `json:"message"`
	Type     string `json:"type"`
}

func HandleConnectionsNotif(w http.ResponseWriter, r *http.Request) {
	SenderID := r.FormValue("userId")
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Register the new client with chatID
	clientsMu.Lock()
	clients[&clientInfo{conn: ws, userID: SenderID}] = true
	clientsMu.Unlock()

	func() {
		for {
			var data MessageData
			// Read messages from WebSocket and decode into MessageData
			err := ws.ReadJSON(&data)

			fmt.Println(data.Type)

			if data.Type == "event" {
				fmt.Println(data.Message)
				dbConnection := database.DB
				groupID := data.UserID
				userID := data.SenderID

				members, err := groupHandlers.GetAllGroupMembersFromDatabase(dbConnection, groupID, userID)
				if err != nil {
					log.Printf("Error fetching group members from database: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}

				for _, member := range members.Members {
					// Отправляем уведомление каждому участнику
					notification := SendMessageResponse{
						Text:   data.Message,
						UserID: member.UserID,
					}

					// Отправляем уведомление в канал broadcast
					broadcast <- notification
				}
			}

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

			// Log the notification data before sending
			fmt.Printf("Sending notification to user %s: %s\n", data.UserID, data.Message)

			response := SendMessageResponse{
				Text:   data.Message,
				UserID: data.UserID,
			}

			// Send the message to broadcast channel
			broadcast <- response
		}
	}()
}

func HandleNotifications() {
	for {
		// Getting a message from the broadcast channel
		msg := <-broadcast

		// Log the notification message before sending

		clientsMu.Lock()
		for client := range clients {
			fmt.Println("clienti:", client.userID)
			fmt.Println("msguserid:", msg.UserID)

			// Send the message only to the specified client
			if client.userID == msg.UserID {
				err := client.conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error writing to client: %v", err)
					client.conn.Close()
					delete(clients, client)
				}
			}
		}
		clientsMu.Unlock()
	}
}
