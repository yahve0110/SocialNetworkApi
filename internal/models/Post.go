package models

type Post struct{
	PostID string `json:"post_id"`
	AuthorID string `json:"author_id"`
	AuthorFirstName string `json:"author_first_name"`
	AuthorLastName string `json:"author_last_name"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	LikesCount int `json:"likes_count"`
	Image string `json:"image"`
	Private string `json:"privacy"`
	AuthorNickname string `json:"author_nickname"`
	PrivateUsersArr []string `json:"private_users"`
}


