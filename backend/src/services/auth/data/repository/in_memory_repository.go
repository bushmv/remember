package repository

import (
	"errors"

	"github.com/bushmv/remember/src/services/auth/data/dto"
)

type User struct {
	id           int64
	username     string
	name         string
	passwordHash string
}

type InMemoryRepository struct {
	table []User
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		table: make([]User, 0),
	}
}

func (r *InMemoryRepository) CreateUser(user *dto.NewUserRequest) *dto.UserCreatedResponse {
	row := User{
		id:           int64(len(r.table)),
		name:         user.Name,
		username:     user.Username,
		passwordHash: user.PasswordHash,
	}
	r.table = append(r.table, row)
	createdUser := &dto.UserCreatedResponse{
		ID:       row.id,
		Name:     row.name,
		Username: row.username,
	}
	return createdUser
}

func (r *InMemoryRepository) FindUserByUsername(username string) (*dto.User, error) {
	for _, row := range r.table {
		if username == row.username {
			user := &dto.User{
				ID:           row.id,
				Name:         row.name,
				Username:     row.username,
				PasswordHash: row.passwordHash,
			}
			return user, nil
		}
	}
	return nil, errors.New("User not found")
}
