package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-todo-api/app"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinanceLifecycle(t *testing.T) {
	SetupTestDB(t)

	router := gin.Default()
	app.SetupRouter(router)

	// 1. Authenticate user
	token, _ := RegisterAndLogin(t, router)

	// ==========================================
	// PART A: Transaction Categories CRUD
	// ==========================================

	// Create Category
	catReq := map[string]interface{}{
		"name": "Food & Beverage",
		"type": "expense",
	}
	w := JSONRequest(router, http.MethodPost, "/transaction-categories", catReq, token)
	assert.Equal(t, http.StatusCreated, w.Code)

	var catCreateRes struct {
		Data struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &catCreateRes)
	require.NoError(t, err)
	categoryID := catCreateRes.Data.ID
	assert.Equal(t, "Food & Beverage", catCreateRes.Data.Name)

	// Update Category
	catUpdateName := "Food & Dining"
	catUpdateReq := map[string]interface{}{
		"name": &catUpdateName,
	}
	w = JSONRequest(router, http.MethodPatch, fmt.Sprintf("/transaction-categories/%d", categoryID), catUpdateReq, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// Get Categories
	w = JSONRequest(router, http.MethodGet, "/transaction-categories", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)
	var catGetRes struct {
		Data []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &catGetRes)
	require.NoError(t, err)
	assert.NotEmpty(t, catGetRes.Data)
	assert.Equal(t, "Food & Dining", catGetRes.Data[0].Name)

	// ==========================================
	// PART B: Transactions CRUD
	// ==========================================

	// Create Transaction
	txDate := time.Now().Format("2006-01-02")
	txReq := map[string]interface{}{
		"title":            "Lunch at Restaurant",
		"amount":           50000.0,
		"type":             "expense",
		"category_id":      categoryID,
		"notes":            "Had lunch with friend",
		"transaction_date": txDate,
	}
	w = JSONRequest(router, http.MethodPost, "/transactions", txReq, token)
	assert.Equal(t, http.StatusCreated, w.Code)

	var txCreateRes struct {
		Data struct {
			ID    uint    `json:"id"`
			Title string  `json:"title"`
			Price float64 `json:"amount"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &txCreateRes)
	require.NoError(t, err)
	transactionID := txCreateRes.Data.ID
	assert.Equal(t, "Lunch at Restaurant", txCreateRes.Data.Title)

	// Update Transaction
	txUpdateTitle := "Lunch at Restaurant (Updated)"
	txUpdateReq := map[string]interface{}{
		"title": &txUpdateTitle,
	}
	w = JSONRequest(router, http.MethodPatch, fmt.Sprintf("/transactions/%d", transactionID), txUpdateReq, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// Get Transactions
	w = JSONRequest(router, http.MethodGet, "/transactions", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// ==========================================
	// PART C: Budgets CRUD
	// ==========================================

	// Create Budget
	currentMonth := time.Now().Format("2006-01")
	budgetReq := map[string]interface{}{
		"category_id": categoryID,
		"amount":      1000000.0,
		"month":       currentMonth,
	}
	w = JSONRequest(router, http.MethodPost, "/budgets", budgetReq, token)
	assert.Equal(t, http.StatusCreated, w.Code)

	var budgetCreateRes struct {
		Data struct {
			ID     uint    `json:"id"`
			Amount float64 `json:"amount"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &budgetCreateRes)
	require.NoError(t, err)
	budgetID := budgetCreateRes.Data.ID
	assert.Equal(t, 1000000.0, budgetCreateRes.Data.Amount)

	// Update Budget
	budgetUpdateAmount := 1500000.0
	budgetUpdateReq := map[string]interface{}{
		"amount": &budgetUpdateAmount,
	}
	w = JSONRequest(router, http.MethodPatch, fmt.Sprintf("/budgets/%d", budgetID), budgetUpdateReq, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// Get Budgets
	w = JSONRequest(router, http.MethodGet, fmt.Sprintf("/budgets?month=%s", currentMonth), nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// ==========================================
	// PART D: Recurring Transactions CRUD
	// ==========================================

	// Create Recurring Transaction
	recReq := map[string]interface{}{
		"title":       "Monthly Gym Subscription",
		"amount":      300000.0,
		"type":        "expense",
		"category_id": categoryID,
		"notes":       "Gym membership fees",
		"frequency":   "monthly",
		"next_run_at": txDate,
	}
	w = JSONRequest(router, http.MethodPost, "/recurring-transactions", recReq, token)
	assert.Equal(t, http.StatusCreated, w.Code)

	var recCreateRes struct {
		Data struct {
			ID    uint   `json:"id"`
			Title string `json:"title"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &recCreateRes)
	require.NoError(t, err)
	recurringID := recCreateRes.Data.ID
	assert.Equal(t, "Monthly Gym Subscription", recCreateRes.Data.Title)

	// Update Recurring Transaction
	recUpdateTitle := "Gym Membership"
	recUpdateReq := map[string]interface{}{
		"title": &recUpdateTitle,
	}
	w = JSONRequest(router, http.MethodPatch, fmt.Sprintf("/recurring-transactions/%d", recurringID), recUpdateReq, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// Get Recurring Transactions
	w = JSONRequest(router, http.MethodGet, "/recurring-transactions", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// ==========================================
	// PART E: Dashboard & Analytics Overview
	// ==========================================

	w = JSONRequest(router, http.MethodGet, "/dashboard/summary", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	w = JSONRequest(router, http.MethodGet, "/finance/summary", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	w = JSONRequest(router, http.MethodGet, fmt.Sprintf("/finance/analytics?month=%s", currentMonth), nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// ==========================================
	// CLEANUP: Deletion
	// ==========================================

	// Delete Recurring
	w = JSONRequest(router, http.MethodDelete, fmt.Sprintf("/recurring-transactions/%d", recurringID), nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// Delete Budget
	w = JSONRequest(router, http.MethodDelete, fmt.Sprintf("/budgets/%d", budgetID), nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// Delete Transaction
	w = JSONRequest(router, http.MethodDelete, fmt.Sprintf("/transactions/%d", transactionID), nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	// Delete Category
	w = JSONRequest(router, http.MethodDelete, fmt.Sprintf("/transaction-categories/%d", categoryID), nil, token)
	assert.Equal(t, http.StatusOK, w.Code)
}
