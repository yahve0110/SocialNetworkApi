package models

type Friendship struct {
    FriendshipID int    `json:"friendship_id"`
    User1ID      int    `json:"user1_id"`
    User2ID      int    `json:"user2_id"`
    Status       string `json:"status"`
    ActionUserID int    `json:"action_user_id"`
}
