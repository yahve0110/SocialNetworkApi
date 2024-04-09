package models

import "time"

type Group struct {
	GroupID          string
	CreatorID        string
	GroupName        string `json:"group_name"`
	GroupDescription string `json:"group_description"`
	GroupImage       string `json:"group_image"`
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

// GroupEvent represents an event in a group
type GroupEvent struct {
	EventID     string    `json:"event_id"`
	GroupID     string    `json:"group_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DateTime    time.Time `json:"date_time"`
	Options     struct {
		Going    []string `json:"going"`
		NotGoing []string `json:"not_going"`
	} `json:"options"`
}

// GroupRequest represents a request to join a group
type GroupRequest struct {
	RequestID string `json:"request_id"`
	GroupID   string `json:"group_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
}

type GroupPost struct {
	PostID    string    `json:"post_id"`
	GroupID   string    `json:"group_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	LikeCount int       `json:"like_count"`
}

type GroupPostComment struct {
	GroupID   string    `json:"group_id"`
	CommentID string    `json:"comment_id"`
	PostID    string    `json:"post_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	LikeCount int       `json:"like_count"`
}

// GroupPostLike represents a like on a group post
type GroupPostLike struct {
	LikeID    string    `json:"like_id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupCommentLike struct {
	LikeID    string `json:"like_id"`
	CommentID string `json:"comment_id"`
	UserID    string `json:"user_id"`
}
