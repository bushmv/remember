package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	jwtManager *JWTManager
}

func NewAuthMiddleware(jwtManager *JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "unauthorized",
				"description": "authorization header is required",
			})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "unauthorized",
				"description": "invalid authorization header format. Use: Bearer <token>",
			})
			return
		}
		tokenString := parts[1]
		token, err := m.jwtManager.VerifyToken(tokenString)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":       "token expired",
					"description": "token has expired",
				})
			case errors.Is(err, jwt.ErrTokenMalformed):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":       "invalid token",
					"description": "token is malformed",
				})
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":       "internal error",
					"description": "failed to verify token",
				})
			}
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":       "invalid token",
				"description": "token is invalid",
			})
			return
		}
		if claims, ok := token.Claims.(*UserClaims); ok {
			c.Set("username", claims.Username)
		}
		c.Next()
	}
}
