package tests

import (
	"testing"
	"time"

	"go-todo-api/config"
	"go-todo-api/model"
	"go-todo-api/scheduler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBudgetAlertScheduler(t *testing.T) {
	SetupTestDB(t)

	// 1. Create a test user
	user := model.User{
		Name:         "Tester User",
		Username:     "tester_user",
		Email:        "tester@example.com",
		PasswordHash: "somehash",
	}
	err := config.DB.Create(&user).Error
	require.NoError(t, err)

	// 2. Create a transaction category (type: expense)
	category := model.TransactionCategory{
		UserID: user.ID,
		Name:   "Food",
		Type:   "expense",
	}
	err = config.DB.Create(&category).Error
	require.NoError(t, err)

	// 3. Create a budget for this month
	currentMonth := time.Now().Format("2006-01")
	budget := model.Budget{
		UserID:     user.ID,
		CategoryID: category.ID,
		Amount:     1000000, // 1 Million
		Month:      currentMonth,
	}
	err = config.DB.Create(&budget).Error
	require.NoError(t, err)

	// Instantiate the alert scheduler
	alertSched := scheduler.NewAlertScheduler()

	// Run scheduler with 0 expenses (0%) -> No alerts should be created
	alertSched.ProcessBudgetAlerts()

	var alerts []model.Alert
	err = config.DB.Where("user_id = ?", user.ID).Find(&alerts).Error
	require.NoError(t, err)
	assert.Len(t, alerts, 0)

	// 4. Insert transaction consuming 85% of budget (850,000)
	tx1 := model.Transaction{
		UserID:          user.ID,
		Title:           "Weekly Groceries",
		Amount:          850000,
		Type:            "expense",
		CategoryID:      category.ID,
		TransactionDate: time.Now(),
	}
	err = config.DB.Create(&tx1).Error
	require.NoError(t, err)

	// Run scheduler -> Should create a budget_warning alert
	alertSched.ProcessBudgetAlerts()

	err = config.DB.Where("user_id = ?", user.ID).Order("created_at asc").Find(&alerts).Error
	require.NoError(t, err)
	require.Len(t, alerts, 1)
	assert.Equal(t, "budget_warning", alerts[0].Type)
	assert.Contains(t, alerts[0].Title, "Peringatan Anggaran")
	assert.Contains(t, alerts[0].Message, "85.0%")

	// Run scheduler again -> Should not duplicate the warning alert
	alertSched.ProcessBudgetAlerts()
	err = config.DB.Where("user_id = ?", user.ID).Find(&alerts).Error
	require.NoError(t, err)
	assert.Len(t, alerts, 1)

	// 5. Insert transaction pushing budget to 105% total (extra 200,000)
	tx2 := model.Transaction{
		UserID:          user.ID,
		Title:           "Dinner Out",
		Amount:          200000,
		Type:            "expense",
		CategoryID:      category.ID,
		TransactionDate: time.Now(),
	}
	err = config.DB.Create(&tx2).Error
	require.NoError(t, err)

	// Run scheduler -> Should create a budget_exceeded alert
	alertSched.ProcessBudgetAlerts()

	err = config.DB.Where("user_id = ?", user.ID).Order("created_at asc").Find(&alerts).Error
	require.NoError(t, err)
	require.Len(t, alerts, 2)
	assert.Equal(t, "budget_warning", alerts[0].Type)
	assert.Equal(t, "budget_exceeded", alerts[1].Type)
	assert.Contains(t, alerts[1].Title, "Anggaran Food Terlampaui!")
	assert.Contains(t, alerts[1].Message, "105.0%")

	// Run scheduler again -> Should not duplicate the exceeded alert
	alertSched.ProcessBudgetAlerts()
	err = config.DB.Where("user_id = ?", user.ID).Find(&alerts).Error
	require.NoError(t, err)
	assert.Len(t, alerts, 2)
}
