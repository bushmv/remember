package dto

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserId   int64
	Username string
	jwt.RegisteredClaims
}
