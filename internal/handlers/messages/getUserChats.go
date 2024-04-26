package messageHandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

type ChatInfo struct {
	ChatID          string `json:"chat_id"`
	InterlocutorID  string `json:"interlocutor_id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	ProfilePicture  string `json:"profile_picture"`
	LastMessage     string `json:"last_message"`
	LastMessageTime string `json:"last_message_time"`
}

func GetUserChats(w http.ResponseWriter, r *http.Request) {
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

	// Запрос к базе данных для получения всех приватных чатов текущего пользователя
	rows, err := dbConnection.Query(`
		SELECT pc.chat_id, CASE WHEN pc.user1_id = ? THEN pc.user2_id ELSE pc.user1_id END AS interlocutor_id,
		u.first_name, u.last_name, u.profile_picture,
		pm.content AS last_message, pm.timestamp AS last_message_time
		FROM privatechat pc
		INNER JOIN users u ON CASE WHEN pc.user1_id = ? THEN pc.user2_id ELSE pc.user1_id END = u.user_id
		LEFT JOIN (
			SELECT chat_id, content, timestamp
			FROM privatechat_messages pm
			WHERE (chat_id, timestamp) IN (
				SELECT chat_id, MAX(timestamp)
				FROM privatechat_messages
				GROUP BY chat_id
			)
		) pm ON pc.chat_id = pm.chat_id
		WHERE pc.user1_id = ? OR pc.user2_id = ?
	`, userID, userID, userID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting user's private chats: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Срез для хранения информации о чатах
	var chatsInfo []ChatInfo

	// Итерация по результатам запроса
	// Итерация по результатам запроса
	for rows.Next() {
		var chatInfo ChatInfo
		var lastMessage, lastMessageTime sql.NullString // Используйте тип sql.NullString для сканирования NULL значений из базы данных

		err := rows.Scan(
			&chatInfo.ChatID, &chatInfo.InterlocutorID,
			&chatInfo.FirstName, &chatInfo.LastName, &chatInfo.ProfilePicture,
			&lastMessage, &lastMessageTime, // Используйте sql.NullString вместо string
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning chat info: %v", err), http.StatusInternalServerError)
			return
		}

		// Присваиваем значения полям ChatInfo, обрабатывая NULL значения
		if lastMessage.Valid {
			chatInfo.LastMessage = lastMessage.String // Если значение не NULL, присваиваем его
		} else {
			chatInfo.LastMessage = "" // Если значение NULL, присваиваем пустую строку
		}

		if lastMessageTime.Valid {
			chatInfo.LastMessageTime = lastMessageTime.String // Если значение не NULL, присваиваем его
		} else {
			chatInfo.LastMessageTime = "" // Если значение NULL, присваиваем пустую строку
		}

		chatsInfo = append(chatsInfo, chatInfo)
	}

	// Возвращение списка информации о чатах в формате JSON
	json.NewEncoder(w).Encode(chatsInfo)
}
