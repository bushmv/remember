package dto

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	PasswordHash string `json:"password hash"`
}
