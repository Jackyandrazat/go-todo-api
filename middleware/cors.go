package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := map[string]bool{
		"http://localhost:3000": true, // web frontend
		"http://localhost:5173": true, // vite
		"http://localhost:8081": true, // local debug
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			if allowedOrigins[origin] {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		c.Writer.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PATCH, DELETE, OPTIONS",
		)

		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Authorization, Content-Type, X-Request-ID",
		)

		c.Writer.Header().Set(
			"Access-Control-Expose-Headers",
			"X-Request-ID",
		)

		if strings.EqualFold(c.Request.Method, http.MethodOptions) {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
