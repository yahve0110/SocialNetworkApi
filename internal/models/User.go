package models

type User struct {
	UserID string `json:"user_id"`
	Username string `json:"username"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Gender string `json:"gender"`
	BirthDate string `json:"birth_date"`
	ProfilePicture string `json:"profilePicture"`
	Role string `json:"role"`
	About string `json:"about"`
	Privacy string `json:"privacy"`

}