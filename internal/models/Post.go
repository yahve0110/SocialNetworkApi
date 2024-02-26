package models

type Post struct{
	PostID string `json:"post_id"`
	AuthorID string `json:"author_id"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	LikesCount int `json:"likes_count"`
	Image string `json:"image"`
	Private string `json:"privacy"`
	AuthorNickname string `json:"author_nickname"`
}


