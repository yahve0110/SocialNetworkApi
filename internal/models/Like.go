package models

type LikePost struct {
	LikeID  int `json:"like_id"`
    PostID  int `json:"post_id"`
    UserID  int `json:"user_id"`
}

type LikeComment struct {
	LikeID  int `json:"like_id"`
    CommentID  int `json:"comment_id"`
    UserID  int `json:"user_id"`
}