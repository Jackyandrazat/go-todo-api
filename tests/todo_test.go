package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"go-todo-api/app"
	"go-todo-api/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoLifecycle(t *testing.T) {
	SetupTestDB(t)

	router := gin.Default()
	app.SetupRouter(router)

	// 1. Authenticate user
	token, _ := RegisterAndLogin(t, router)

	// 2. Create Todo
	createReq := map[string]interface{}{
		"title":       "Learn Integration Testing in Go",
		"description": "Write integration tests for the API",
		"priority":    "high",
		"category":    "Work",
	}

	w := JSONRequest(router, http.MethodPost, "/todos", createReq, token)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createRes struct {
		Data model.Todo `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &createRes)
	require.NoError(t, err)
	assert.Equal(t, "Learn Integration Testing in Go", createRes.Data.Title)
	todoID := createRes.Data.ID

	// 3. Get Todos
	w = JSONRequest(router, http.MethodGet, "/todos", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var getRes struct {
		Data struct {
			Items []model.Todo `json:"items"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &getRes)
	require.NoError(t, err)
	assert.NotEmpty(t, getRes.Data.Items)
	assert.Equal(t, todoID, getRes.Data.Items[0].ID)

	// 4. Update Todo
	isDone := true
	updateReq := map[string]interface{}{
		"title": "Learn Integration Testing in Go (Updated)",
		"done":  &isDone,
	}

	w = JSONRequest(router, http.MethodPatch, fmt.Sprintf("/todos/%d", todoID), updateReq, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var updateRes struct {
		Data model.Todo `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &updateRes)
	require.NoError(t, err)
	assert.Equal(t, "Learn Integration Testing in Go (Updated)", updateRes.Data.Title)
	assert.True(t, updateRes.Data.Done)

	// 5. Delete Todo
	w = JSONRequest(router, http.MethodDelete, fmt.Sprintf("/todos/%d", todoID), nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify deletion
	w = JSONRequest(router, http.MethodGet, "/todos", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &getRes)
	require.NoError(t, err)
	assert.Empty(t, getRes.Data.Items)
}
