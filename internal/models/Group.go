package models

type Group struct {
	GroupID          int    `json:"group_id"`
	GroupName        string `json:"group_name"`
	GroupDescription string `json:"group_description"`
	CreatorID        int    `json:"creator_id"`
	CreationDate     string `json:"creation_date"`
	Status           string `json:"status"`
}

type GroupMember struct {
	MembershipID int `json:"membership_id"`
	GroupID      int `json:"group_id"`
	MemberID     int `json:"member_id"`
}

type GroupIvent struct {
	EventID int `json:"group_event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"event_datetime "`
}
