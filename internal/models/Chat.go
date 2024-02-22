package models

type Chat struct {
	MessageId int `json:"message_id"`
	ChatID    int    `json:"chat_id"`
    User1ID  int    `json:"user1_id"`
	User2ID  int    `json:"user2_id"`
    Content   string `json:"content"`
    Timestamp string `json:"created_at"`
}

type GroupChat struct {
	ChatID int `json:"chat_id"`
	ChatName string `json:"chat_name"`
	CreatorID int `json:"creator_id"`
	CreationDate string `json:"creation_date"`
}


type GroupChatMember struct {
	MemberID int `json:"member_id"`
	ChatID int `json:"chat_id"`
}

type GroupChatMessage struct {
	MessageID int `json:"message_id"`
	Content string `json:"content"`
	SenderID int `json:"sender_id"`
	ChatID int `json:"chat_id"`
	Timestamp string `json:"created_at"`

}