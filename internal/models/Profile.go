package models

type Profile struct {
	ProfileID int `json:"profile_id"`
	FollowerID int `json:"follower_id"`
	ProfileStatus string `json:"profile_status"`
}