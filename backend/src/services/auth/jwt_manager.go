package auth

import (
	"fmt"

	"github.com/bushmv/remember/src/services/auth/data/dto"
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

func (m *JWTManager) CreateJWTToken(claims dto.UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *JWTManager) keyFunc(t *jwt.Token) (any, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}
	return m.secretKey, nil
}

func (m *JWTManager) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.UserClaims{}, m.keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token not valid: %v", token)
	}
	if _, ok := token.Claims.(*dto.UserClaims); !ok {
		return nil, fmt.Errorf("cannot cast to *UserClaims")
	}
	return token, nil
}
