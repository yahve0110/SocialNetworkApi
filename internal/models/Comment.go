package models


type Comment struct {
	CommentID int `json:"comment_id"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	AuthorID int `json:"author_id"`
	PostID int `json:"post_id"`
	Timestamp string `json:"comment_created_at"`
	Image string `json:"image"`
}

