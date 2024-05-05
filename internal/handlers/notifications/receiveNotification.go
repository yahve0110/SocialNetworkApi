package notifications

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"time"

	"github.com/google/uuid"
)

// Notification represents a notification entity
type Notification struct {
	ID         string `json:"notification_id"`
	ReceiverID string `json:"receiver_id"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	SenderID   string `json:"sender_id"`
	GroupID    string `json:"group_id"`
}

// ReceiveNotification обрабатывает HTTP-запрос для получения уведомлений
func ReceiveNotification(w http.ResponseWriter, r *http.Request) {
	var notification Notification
	dbConnection := database.DB

	// Get the user ID based on the current user's session
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get the user ID based on the current user's session
	userID, err := helpers.GetUserIDFromSession(dbConnection, cookie.Value)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Декодирование JSON-данных из тела запроса в структуру Notification
	err = json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Failed to decode notification data", http.StatusBadRequest)
		return
	}

	if notification.Type == "group_request" && notification.ReceiverID == "" {
		creatorID, err := getGroupCreatorID(notification.GroupID)
		if err != nil {
			http.Error(w, "Failed to get group creator ID: "+err.Error(), http.StatusInternalServerError)
			return
		}

		notification.ReceiverID = creatorID
	}

	notification.ID = uuid.New().String()
	// Установка текущего времени в поле CreatedAt
	notification.CreatedAt = time.Now().Format(time.RFC3339)

	notification.SenderID = userID

	if notification.Type == "group_event" {
		// Получаем список участников группы
		memberIDs, err := getGroupMembers(notification.GroupID)
		if err != nil {
			http.Error(w, "Failed to get group members: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for _, memberID := range memberIDs {
			eventNotification := Notification{
				ID:         uuid.New().String(),
				ReceiverID: memberID,
				Type:       "group_event",
				Content:    notification.Content,
				CreatedAt:  time.Now().Format(time.RFC3339),
				SenderID:   notification.SenderID,
				GroupID:    notification.GroupID,
			}

			err := SaveNotification(eventNotification)
			if err != nil {
				http.Error(w, "Failed to save notification: "+err.Error(), http.StatusInternalServerError)
				return
			}

		}
	}

	// Сохранение уведомления в базу данных
	err = SaveNotification(notification)
	if err != nil {
		http.Error(w, "Failed to save notification", http.StatusInternalServerError)
		return
	}

	// Отправка успешного ответа
	w.WriteHeader(http.StatusOK)
}

func SaveNotification(notification Notification) error {
	db := database.DB

	query := `
		INSERT INTO notifications (notification_id, receiver_id, type, content, created_at, sender_id, group_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query, notification.ID, notification.ReceiverID, notification.Type, notification.Content, notification.CreatedAt, notification.SenderID, notification.GroupID)
	if err != nil {
		log.Printf("Failed to insert notification into database: %v", err)
		return err
	}

	return nil
}

func getGroupCreatorID(groupID string) (string, error) {
	var creatorID string
	db := database.DB

	err := db.QueryRow("SELECT creator_id FROM groups WHERE group_id = ?", groupID).Scan(&creatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	return creatorID, nil
}

func getGroupMembers(groupID string) ([]string, error) {
	var memberIDs []string
	db := database.DB

	rows, err := db.Query("SELECT user_id FROM group_members WHERE group_id = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var memberID string
		err := rows.Scan(&memberID)
		if err != nil {
			return nil, err
		}
		memberIDs = append(memberIDs, memberID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memberIDs, nil
}


