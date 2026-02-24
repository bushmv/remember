package dto

type UserCreatedResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}
