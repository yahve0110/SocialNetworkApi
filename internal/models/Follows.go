package models


type Follows struct {
	FollowsID int `json:"follow_id"`
	UserFollowedStatus string `json:"user_followed_status"`
	UserFollowed string `json:"user_followed"`
	UserFollowing string `json:"user_following"`
}