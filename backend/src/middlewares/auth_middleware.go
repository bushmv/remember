package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/bushmv/remember/src/services/auth"
	"github.com/bushmv/remember/src/services/auth/data/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format. Use: Bearer <token>",
			})
			return
		}
		tokenString := parts[1]
		token, err := m.jwtManager.VerifyToken(tokenString)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "token has expired",
				})
			case errors.Is(err, jwt.ErrTokenMalformed):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "token is malformed",
				})
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "token is invalid",
			})
			return
		}
		if claims, ok := token.Claims.(*dto.UserClaims); ok {
			c.Set("userId", claims.ID)
			c.Set("username", claims.Username)
		}
		c.Next()
	}
}
