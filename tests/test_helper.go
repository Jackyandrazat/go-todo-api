package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"go-todo-api/config"
	"go-todo-api/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func SetupTestDB(t *testing.T) {
	os.Setenv("APP_ENV", "test")

	config.LoadConfig()
	config.ConnectTestDB()

	err := config.DB.Migrator().DropTable(
		&model.UserSession{},
		&model.Alert{},
		&model.RecurringTransaction{},
		&model.Budget{},
		&model.Transaction{},
		&model.TransactionCategory{},
		&model.Note{},
		&model.Todo{},
		&model.User{},
	)
	require.NoError(t, err)

	err = config.DB.AutoMigrate(
		&model.User{},
		&model.Todo{},
		&model.Note{},
		&model.TransactionCategory{},
		&model.Transaction{},
		&model.Budget{},
		&model.RecurringTransaction{},
		&model.Alert{},
		&model.UserSession{},
	)
	require.NoError(t, err)
}

func JSONRequest(
	router *gin.Engine,
	method string,
	url string,
	body interface{},
	token string,
) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(
		method,
		url,
		bytes.NewBuffer(reqBody),
	)

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set(
			"Authorization",
			"Bearer "+token,
		)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func RegisterAndLogin(t *testing.T, router *gin.Engine) (string, uint) {
	register := map[string]interface{}{
		"name":     "Test User",
		"username": "tester",
		"email":    "test@example.com",
		"password": "password123",
	}

	w := JSONRequest(router, "POST", "/auth/register", register, "")
	require.Equal(t, 201, w.Code)

	login := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}

	w = JSONRequest(router, "POST", "/auth/login", login, "")
	require.Equal(t, 200, w.Code)

	var res struct {
		Data struct {
			AccessToken string `json:"access_token"`
			User        struct {
				ID uint `json:"id"`
			} `json:"user"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)

	return res.Data.AccessToken, res.Data.User.ID
}
