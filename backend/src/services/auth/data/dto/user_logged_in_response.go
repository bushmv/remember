package dto

type UserLoggedInResponse struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
