package main

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type JWTManager struct {
	secretKey []byte
}

func NewJWTManager(secretKey []byte) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
	}
}

func (m *JWTManager) CreateJWTToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *JWTManager) VerifyToken(tokenString string) (bool, error) {
	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	})
	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, fmt.Errorf("token not valid: %v", token)
	}
	return true, nil
}
