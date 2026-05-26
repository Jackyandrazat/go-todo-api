package middleware

import (
	"net/http"
	"strings"

	"go-todo-api/response"
	"go-todo-api/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Error(
				c,
				http.StatusUnauthorized,
				"authorization header required",
				nil,
			)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(
				c,
				http.StatusUnauthorized,
				"invalid authorization format",
				nil,
			)
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			response.Error(
				c,
				http.StatusUnauthorized,
				"invalid or expired token",
				nil,
			)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
