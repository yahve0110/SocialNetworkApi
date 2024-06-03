package notifications

import (
	"encoding/json"
	"net/http"

	database "social/internal/db"
	"social/internal/helpers"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {

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

	notifications, err := fetchNotifications(userID)
	if err != nil {
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func fetchNotifications(userID string) ([]Notification, error) {
	db := database.DB

	query := `
		SELECT notification_id, receiver_id, type, content, created_at, sender_id, group_id
		FROM notifications
		WHERE receiver_id = ?
		ORDER BY created_at DESC
		LIMIT 15
	`

	// Выполнение SQL-запроса
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		if err := rows.Scan(&notification.ID, &notification.ReceiverID, &notification.Type, &notification.Content, &notification.CreatedAt, &notification.SenderID, &notification.GroupID); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
