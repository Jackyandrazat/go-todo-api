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

func TestNoteLifecycle(t *testing.T) {
	SetupTestDB(t)

	router := gin.Default()
	app.SetupRouter(router)

	// 1. Authenticate user
	token, _ := RegisterAndLogin(t, router)

	// 2. Create Note
	createReq := map[string]interface{}{
		"title":    "My Integration Test Note",
		"content":  "Testing notes endpoint behaves correctly",
		"category": "Personal",
	}

	w := JSONRequest(router, http.MethodPost, "/notes", createReq, token)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createRes struct {
		Data model.Note `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &createRes)
	require.NoError(t, err)
	assert.Equal(t, "My Integration Test Note", createRes.Data.Title)
	assert.False(t, createRes.Data.IsPinned)
	noteID := createRes.Data.ID

	// 3. Get Notes
	w = JSONRequest(router, http.MethodGet, "/notes", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var getRes struct {
		Data struct {
			Items []model.Note `json:"items"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &getRes)
	require.NoError(t, err)
	assert.NotEmpty(t, getRes.Data.Items)
	assert.Equal(t, noteID, getRes.Data.Items[0].ID)

	// 4. Update Note
	updateReq := map[string]interface{}{
		"title": "My Integration Test Note (Updated)",
	}

	w = JSONRequest(router, http.MethodPatch, fmt.Sprintf("/notes/%d", noteID), updateReq, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var updateRes struct {
		Data model.Note `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &updateRes)
	require.NoError(t, err)
	assert.Equal(t, "My Integration Test Note (Updated)", updateRes.Data.Title)

	// 5. Toggle Pin Status
	pinReq := map[string]interface{}{
		"is_pinned": true,
	}
	w = JSONRequest(router, http.MethodPatch, fmt.Sprintf("/notes/%d/pin", noteID), pinReq, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var pinRes struct {
		Data model.Note `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &pinRes)
	require.NoError(t, err)
	assert.True(t, pinRes.Data.IsPinned)

	// 6. Delete Note
	w = JSONRequest(router, http.MethodDelete, fmt.Sprintf("/notes/%d", noteID), nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify deletion
	w = JSONRequest(router, http.MethodGet, "/notes", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &getRes)
	require.NoError(t, err)
	assert.Empty(t, getRes.Data.Items)
}
