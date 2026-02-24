package repository

import "github.com/bushmv/remember/src/services/auth/data/dto"

type UserRepository interface {
	CreateUser(user *dto.NewUserRequest) *dto.UserCreatedResponse
	FindUserByUsername(username string) (*dto.User, error)
}
