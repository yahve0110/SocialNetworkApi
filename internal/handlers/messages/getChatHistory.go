package messageHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social/internal/db"
	"social/internal/helpers"
	"time"
)

// ChatMessage represents a chat message
type ChatMessage struct {
	MessageID       string    `json:"message_id"`
	ChatID          string    `json:"chat_id"`
	MessageAuthorID string    `json:"message_author_id"`
	Content         string    `json:"content"`
	Timestamp       time.Time `json:"timestamp"`
}

// MessageWithUser represents a message with user information
type MessageWithUser struct {
	ChatID          string    `json:"chat_id"`
	MessageAuthorID string    `json:"message_author_id"`
	Content         string    `json:"content"`
	Timestamp       time.Time `json:"timestamp"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	ProfilePicture  string    `json:"profile_picture"`
}

// ChatInfo represents information about a chat
type UserChatInfo struct {
	ChatID         string `json:"chat_id"`
	InterlocutorID string `json:"interlocutor_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
}

// GetChatHistory retrieves all messages in a chat along with user information
func GetChatHistory(w http.ResponseWriter, r *http.Request) {
	dbConnection := database.DB

	// Get the user ID based on the current user's session
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

	// Get the chat ID from the request URL or request body
	chatID := r.FormValue("chat_id")
	fmt.Println("chat id: ", chatID)

	// Check if the current user is a participant in the chat
	var exists bool
	err = dbConnection.QueryRow("SELECT EXISTS(SELECT 1 FROM privatechat WHERE chat_id = ? AND (user1_id = ? OR user2_id = ?))", chatID, userID, userID).Scan(&exists)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking chat existence: %v", err), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Chat not found or user is not a participant", http.StatusNotFound)
		return
	}

	// Retrieve information about the chat's interlocutor
	var chatInfo UserChatInfo
	err = dbConnection.QueryRow(`
		SELECT CASE WHEN user1_id = ? THEN user2_id ELSE user1_id END AS interlocutor_id, u.first_name, u.last_name, u.profile_picture
		FROM privatechat pc
		INNER JOIN users u ON CASE WHEN pc.user1_id = ? THEN pc.user2_id ELSE pc.user1_id END = u.user_id
		WHERE pc.chat_id = ?
	`, userID, userID, chatID).Scan(&chatInfo.InterlocutorID, &chatInfo.FirstName, &chatInfo.LastName, &chatInfo.ProfilePicture)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving chat information: %v", err), http.StatusInternalServerError)
		return
	}
	chatInfo.ChatID = chatID

	// Retrieve all messages in the chat along with user information
	rows, err := dbConnection.Query(`
		SELECT m.message_author_id, m.content, m.timestamp, u.first_name, u.last_name, u.profile_picture
		FROM privatechat_messages AS m
		INNER JOIN users AS u ON m.message_author_id = u.user_id
		WHERE m.chat_id = ?
		ORDER BY m.timestamp ASC
	`, chatID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving messages: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Slice to store messages with user information
	var messagesWithUsers []MessageWithUser

	// Iterate over the rows and populate the slice
	for rows.Next() {
		var messageWithUser MessageWithUser
		err := rows.Scan(&messageWithUser.MessageAuthorID, &messageWithUser.Content, &messageWithUser.Timestamp, &messageWithUser.FirstName, &messageWithUser.LastName, &messageWithUser.ProfilePicture)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning message: %v", err), http.StatusInternalServerError)
			return
		}
		messagesWithUsers = append(messagesWithUsers, messageWithUser)
	}

	// Combine chat information and messages with user information
	data := struct {
		UserChatInfo         `json:"chat_info"`
		MessagesWithUser []MessageWithUser `json:"messages_with_user"`
	}{
		chatInfo,
		messagesWithUsers,
	}

	// Encode the data as JSON and send the response
	json.NewEncoder(w).Encode(data)
}
