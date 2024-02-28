package models

import "time"

type Group struct {
	GroupID          string
	CreatorID        string
	GroupName        string    `json:"group_name"`
	GroupDescription string    `json:"group_description"`
	GroupImage       string    `json:"group_image"`
	CreationDate     time.Time

}

type GroupInvitation struct {
    GroupID    string `json:"group_id"`
    InviterID  string `json:"inviter_id"`
    ReceiverID string `json:"receiver_id"`
    Status     string `json:"status"`
}

type GroupMember struct {
	MembershipID int `json:"membership_id"`
	GroupID      int `json:"group_id"`
	MemberID     int `json:"member_id"`
}

type GroupEvent struct {
	EventID     int    `json:"group_event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"event_datetime "`
}

// GroupRequest represents a request to join a group
type GroupRequest struct {
	RequestID string `json:"request_id"`
	GroupID   string `json:"group_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
}
