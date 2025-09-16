package entities

type User struct {
	ID             uint     `json:"id"`
	Role           Role     `json:"role"`
	Nickname       string   `json:"nickname"`
	Password       string   `json:"password"`
	ProfilePicture string   `json:"pfp"`
}
