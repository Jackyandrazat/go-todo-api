package middleware

import (
	"net/http"

	"go-todo-api/response"
	"go-todo-api/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				utils.Logger.Error(
					"panic recovered",
					zap.Any("panic", rec),
				)

				response.Error(
					c,
					http.StatusInternalServerError,
					"internal server error",
					nil,
				)

				c.Abort()
			}
		}()

		c.Next()
	}
}
