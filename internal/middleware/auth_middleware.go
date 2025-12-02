package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"ethos/pkg/jwt"
)

// AuthMiddleware validates JWT access tokens and injects user ID into context
func AuthMiddleware(tokenGen *jwt.TokenGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
				"code":  "AUTH_TOKEN_INVALID",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
				"code":  "AUTH_TOKEN_INVALID",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		userID, err := tokenGen.ValidateAccessToken(token)
		if err != nil {
			if err.Error() == "token expired" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token has expired",
					"code":  "AUTH_TOKEN_EXPIRED",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid token",
					"code":  "AUTH_TOKEN_INVALID",
				})
			}
			c.Abort()
			return
		}

		// Inject user ID into context
		c.Set("user_id", userID)
		c.Next()
	}
}

