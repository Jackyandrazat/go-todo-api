package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		requestIDRaw, _ := c.Get("request_id")
		requestID, _ := requestIDRaw.(string)

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf(
			`{"level":"info","request_id":"%s","method":"%s","path":"%s","status":%d,"duration_ms":%d,"client_ip":"%s"}`,
			requestID,
			method,
			path,
			statusCode,
			duration.Milliseconds(),
			clientIP,
		)
	}
}
