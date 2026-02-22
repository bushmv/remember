package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtManager JWTManager
	db         *Database
}

func NewAuthService(db *Database, jwtManager JWTManager) *AuthService {
	return &AuthService{
		jwtManager: jwtManager,
		db:         db,
	}
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (as *AuthService) Register(c *gin.Context) {
	var ur UserRequest
	if err := c.ShouldBindJSON(&ur); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	passwordHash, err := as.hashPassword(ur.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot create new user",
		})
	}
	newUser := &User{
		ID:           0,
		Username:     ur.Username,
		Name:         "",
		PasswordHash: passwordHash,
	}
	as.db.users[ur.Username] = newUser
	c.JSON(http.StatusOK, gin.H{
		"status":   "created",
		"username": ur.Username,
	})
}

func (as *AuthService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (as *AuthService) Login(c *gin.Context) {
	var ur UserRequest
	if err := c.ShouldBindJSON(&ur); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var user *User
	user, exists := as.db.users[ur.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":      "unauthorized",
			"description": "no username: " + ur.Username,
		})
		return
	}
	if as.WrongPassword(ur.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":       "authorization error",
			"description": "wrong login or/and password",
		})
		return
	}

	claims := UserClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	tokenString, err := as.jwtManager.CreateJWTToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"message":  "Authentication successful",
		"username": user.Username,
		"token":    tokenString,
	})
}

func (as *AuthService) WrongPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err != nil
}

func (as *AuthService) RegisterRoutes(r *gin.Engine) {
	r.POST("/register", as.Register)
	r.POST("/login", as.Login)
}
