package tests

import (
	"net/http"
	"testing"

	"go-todo-api/app"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	SetupTestDB(t)
	router := gin.Default()
	app.SetupRouter(router)

	register := map[string]interface{}{
		"name":     "Test User",
		"username": "tester",
		"email":    "test@example.com",
		"password": "password123",
	}

	w := JSONRequest(
		router,
		http.MethodPost,
		"/auth/register",
		register,
		"",
	)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLogin(t *testing.T) {
	SetupTestDB(t)
	router := gin.Default()
	app.SetupRouter(router)

	register := map[string]interface{}{
		"name":     "Test User",
		"username": "tester",
		"email":    "test@example.com",
		"password": "password123",
	}

	JSONRequest(
		router,
		http.MethodPost,
		"/auth/register",
		register,
		"",
	)

	login := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}

	w := JSONRequest(
		router,
		http.MethodPost,
		"/auth/login",
		login,
		"",
	)

	assert.Equal(t, http.StatusOK, w.Code)
}
