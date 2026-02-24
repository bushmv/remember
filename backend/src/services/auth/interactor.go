package auth

import (
	"github.com/bushmv/remember/src/services/auth/data/dto"
	"github.com/bushmv/remember/src/services/auth/data/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthInteractor struct {
	repo       repository.UserRepository
	jwtManager JWTManager
}

func NewAuthInteractor(repo repository.UserRepository, jwtManager JWTManager) *AuthInteractor {
	return &AuthInteractor{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (i *AuthInteractor) CreateUser(userRequest *dto.NewUserRequest) *dto.UserCreatedResponse {
	return i.repo.CreateUser(userRequest)
}
func (i *AuthInteractor) User(username string) (*dto.User, error) {
	return i.repo.FindUserByUsername(username)
}

func (i *AuthInteractor) CreateJWTToken(claims dto.UserClaims) (string, error) {
	return i.jwtManager.CreateJWTToken(claims)
}

func (i *AuthInteractor) VerifyToken(tokenString string) (*jwt.Token, error) {
	return i.jwtManager.VerifyToken(tokenString)
}

func (i *AuthInteractor) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (i *AuthInteractor) CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
