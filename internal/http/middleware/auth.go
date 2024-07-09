package middleware

import (
	"strings"

	"github.com/CnTeng/rx-todo/internal/database"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/token/add" || c.Request.URL.Path == "/user/add" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
		}
		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Authorization header format must be Bearer {token}"})
		}

		token := authParts[1]

		user, err := db.AuthToken(&token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
		}
		c.Set("user", user)

		c.Next()
	}
}
