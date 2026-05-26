package repository

import (
	"go-todo-api/config"
	"go-todo-api/dto"
	"go-todo-api/model"
	"time"
)

type FinanceAnalyticsRepository struct{}

func NewFinanceAnalyticsRepository() *FinanceAnalyticsRepository {
	return &FinanceAnalyticsRepository{}
}

func (r *FinanceAnalyticsRepository) GetCategoryBreakdown(
	userID uint,
	txType string,
	startDate time.Time,
	endDate time.Time,
) ([]dto.CategoryAnalytics, error) {
	var result []dto.CategoryAnalytics

	err := config.DB.
		Model(&model.Transaction{}).
		Joins("JOIN transaction_categories tc ON tc.id = transactions.category_id").
		Select("tc.name as category, SUM(transactions.amount) as amount").
		Where(
			"transactions.user_id = ? AND transactions.type = ? AND transactions.transaction_date BETWEEN ? AND ?",
			userID,
			txType,
			startDate,
			endDate,
		).
		Group("tc.name").
		Order("amount DESC").
		Scan(&result).
		Error

	return result, err
}

type dailyRaw struct {
	Date   time.Time
	Type   string
	Amount float64
}

func (r *FinanceAnalyticsRepository) GetDailyTransactions(
	userID uint,
	startDate time.Time,
	endDate time.Time,
) ([]dailyRaw, error) {
	var rows []dailyRaw

	err := config.DB.
		Model(&model.Transaction{}).
		Select("DATE(transaction_date) as date, type, SUM(amount) as amount").
		Where(
			"user_id = ? AND transaction_date BETWEEN ? AND ?",
			userID,
			startDate,
			endDate,
		).
		Group("DATE(transaction_date), type").
		Order("date ASC").
		Scan(&rows).
		Error

	return rows, err
}
