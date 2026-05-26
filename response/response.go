package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success   bool        `json:"success" example:"true"`
	Message   string      `json:"message" example:"operation successful"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty" example:"req_123456"`
}

type ErrorResponse struct {
	Success   bool        `json:"success" example:"false"`
	Message   string      `json:"message" example:"validation failed"`
	Errors    interface{} `json:"errors,omitempty"`
	RequestID string      `json:"request_id,omitempty" example:"req_123456"`
}

func getRequestID(c *gin.Context) string {
	requestID, exists := c.Get("request_id")
	if !exists {
		return ""
	}

	id, ok := requestID.(string)
	if !ok {
		return ""
	}

	return id
}

func Success(
	c *gin.Context,
	statusCode int,
	message string,
	data interface{},
) {
	c.JSON(statusCode, APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: getRequestID(c),
	})
}

func Error(
	c *gin.Context,
	statusCode int,
	message string,
	errors interface{},
) {
	c.JSON(statusCode, ErrorResponse{
		Success:   false,
		Message:   message,
		Errors:    errors,
		RequestID: getRequestID(c),
	})
}

func Unauthorized(c *gin.Context, message string) {
	Error(
		c,
		http.StatusUnauthorized,
		message,
		nil,
	)
}

func BadRequest(
	c *gin.Context,
	message string,
	errors interface{},
) {
	Error(
		c,
		http.StatusBadRequest,
		message,
		errors,
	)
}

func InternalServerError(c *gin.Context, message string) {
	Error(
		c,
		http.StatusInternalServerError,
		message,
		nil,
	)
}
