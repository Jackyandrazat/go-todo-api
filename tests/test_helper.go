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
