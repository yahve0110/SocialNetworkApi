package models

type Comment struct {
	CommentID string `json:"comment_id"`
	AuthorFirstName string `json:"author_first_name"`
	AuthorLastName string `json:"author_last_name"`
	AuthorAvatar string `json:"author_avatar"`
	Content   string `json:"content"`
	AuthorID  string `json:"author_id"`
	PostID    string   `json:"post_id"`
	AuthorNickname  string `json:"author_nickname"`
	CreatedAt string `json:"comment_created_at"`
	Image     string `json:"image"`
	LikesCount  int `json:"likes_count"`
}
