package main

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey []byte
}

func NewJWTManager(secretKey []byte) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
	}
}

func (m *JWTManager) CreateJWTToken(claims UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *JWTManager) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])

		}
		return m.secretKey, nil
	},
		jwt.WithIssuer("remember"))

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token not valid: %v", token)
	}
	if _, ok := token.Claims.(*UserClaims); !ok {
		return nil, fmt.Errorf("cannot cast tp *UserClaims")
	}
	return token, nil
}
