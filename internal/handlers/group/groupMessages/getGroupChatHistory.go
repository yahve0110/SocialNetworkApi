package groupChat

import (
    "encoding/json"
    "net/http"
    database "social/internal/db"
)

type GroupChatMessage struct {
    MessageID  string `json:"message_id"`
    Content    string `json:"content"`
    AuthorID   string `json:"author_id"`
    ChatID     string `json:"chat_id"`
    CreatedAt  string `json:"created_at"`
    AuthorName string `json:"author_name"`
}

func GetGroupChatHistory(w http.ResponseWriter, r *http.Request) {
    chatID := r.URL.Query().Get("chat_id")

    chatHistory, err := fetchGroupChatHistory(chatID)
    if err != nil {
        http.Error(w, "Failed to fetch group chat history", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(chatHistory)
}

func fetchGroupChatHistory(chatID string) ([]GroupChatMessage, error) {
    db := database.DB

    query := `
        SELECT gcm.message_id, gcm.content, gcm.author_id, gcm.chat_id, gcm.created_at, u.first_name
        FROM group_chat_messages gcm
        INNER JOIN users u ON gcm.author_id = u.user_id
        WHERE gcm.chat_id = ?
        ORDER BY gcm.created_at ASC

        `
    rows, err := db.Query(query, chatID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var chatHistory []GroupChatMessage
    for rows.Next() {
        var message GroupChatMessage
        if err := rows.Scan(&message.MessageID, &message.Content, &message.AuthorID, &message.ChatID, &message.CreatedAt, &message.AuthorName); err != nil {
            return nil, err
        }
        chatHistory = append(chatHistory, message)
    }

    return chatHistory, nil
}
