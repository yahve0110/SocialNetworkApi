package groupChat

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
)

type GroupChatInfo struct {
	ChatID          string         `json:"chat_id"`
	ChatName        string         `json:"chat_name"`
	CreatorID       string         `json:"creator_id"`
	LastMessage     sql.NullString `json:"last_message"`
    LastMessageTime sql.NullString `json:"last_message_time"`
}

func GetGroupChats(w http.ResponseWriter, r *http.Request) {
	dbConnection := database.DB

	// Get the session ID cookie from the request
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

	// Query the database to retrieve all group chats for the current user
	query := `
	SELECT gc.chat_id, gc.chat_name, gc.creator_id,
	gcm.content AS last_message, gcm.created_at AS last_message_time
FROM group_chat gc
INNER JOIN group_chat_members gcmem ON gc.chat_id = gcmem.chat_id
LEFT JOIN (
 SELECT chat_id, content, created_at
 FROM group_chat_messages
 WHERE (chat_id, created_at) IN (
	 SELECT chat_id, MAX(created_at)
	 FROM group_chat_messages
	 GROUP BY chat_id
 )
) gcm ON gc.chat_id = gcm.chat_id
WHERE gcmem.member_id = ?



	`
	rows, err := dbConnection.Query(query, userID)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Error getting user's group chats", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Slice to store information about group chats
	var groupChatsInfo []GroupChatInfo

	// Iterate over the query results and add chat information to the slice
	for rows.Next() {
		var chatInfo GroupChatInfo
		err := rows.Scan(
			&chatInfo.ChatID, &chatInfo.ChatName, &chatInfo.CreatorID,
			&chatInfo.LastMessage, &chatInfo.LastMessageTime,
		)
		if err != nil {
			log.Printf("Error scanning group chat info: %v", err)
			http.Error(w, "Error scanning group chat info", http.StatusInternalServerError)
			return
		}
		groupChatsInfo = append(groupChatsInfo, chatInfo)
	}

	// Log the retrieved group chat information
	log.Printf("Retrieved group chat information: %+v", groupChatsInfo)

	// Return the list of group chat information in JSON format
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groupChatsInfo); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

}
