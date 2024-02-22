package models

type Post struct{
	PostID int `json:"post_id"`
	AuthorID int `json:"author_id"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	LikesCount int `json:"likes_count"`
	Image string `json:"image"`
}


