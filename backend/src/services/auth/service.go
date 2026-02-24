package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bushmv/remember/src/services/auth/data/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	interactor *AuthInteractor
}

func NewAuthService(interactor *AuthInteractor) *AuthService {
	return &AuthService{
		interactor: interactor,
	}
}

func (s *AuthService) SignUp(c *gin.Context) {
	var r dto.SignUpRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	passwordHash, err := s.interactor.hashPassword(r.Password)
	if err != nil {
		err = fmt.Errorf("error when creating a new user")
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	dtoUser := &dto.NewUserRequest{
		Username:     r.Username,
		PasswordHash: passwordHash,
	}
	resp := s.interactor.CreateUser(dtoUser)
	c.JSON(http.StatusCreated, resp)
}

func (s *AuthService) SignIn(c *gin.Context) {
	var r dto.SignInRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	user, err := s.interactor.User(r.Username)
	if err != nil {
		err = fmt.Errorf("Username and/or password not correct")
		c.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}
	if err := s.interactor.CheckPassword(r.Password, user.PasswordHash); err != nil {
		err = fmt.Errorf("Username and/or password not correct")
		c.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}
	claims := dto.UserClaims{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	tokenString, err := s.interactor.CreateJWTToken(claims)
	if err != nil {
		err = errors.New("Internal server error")
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	resp := dto.UserLoggedInResponse{
		Id:       user.ID,
		Username: user.Username,
		Token:    tokenString,
	}
	c.JSON(http.StatusOK, resp)
}

func (as *AuthService) BindRoutes(r *gin.Engine) {
	fmt.Println("RUN")
	r.POST("/api/v1/auth/signup", as.SignUp)
	r.POST("/api/v1/auth/signin", as.SignIn)
}

func ErrorResponse(err error) map[string]any {
	return map[string]any{
		"error": err.Error(),
	}
}
