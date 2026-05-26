package repository

import (
	"time"

	"go-todo-api/config"
	"go-todo-api/model"
)

type TransactionRepository struct{}

type FinanceSummary struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) FindAllByUserID(
	userID uint,
	page int,
	limit int,
	txType *string,
	categoryID *uint,
	startDate *time.Time,
	endDate *time.Time,
) ([]model.Transaction, int64, error) {
	var transactions []model.Transaction
	var total int64

	query := config.DB.
		Model(&model.Transaction{}).
		Preload("Category").
		Where("user_id = ?", userID)

	if txType != nil {
		query = query.Where("type = ?", *txType)
	}

	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if startDate != nil {
		query = query.Where("transaction_date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("transaction_date <= ?", *endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := query.
		Order("transaction_date DESC").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).
		Error

	return transactions, total, err
}

func (r *TransactionRepository) FindByIDAndUserID(
	id uint,
	userID uint,
) (*model.Transaction, error) {
	var transaction model.Transaction

	err := config.DB.
		Preload("Category").
		Where("id = ? AND user_id = ?", id, userID).
		First(&transaction).
		Error

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *TransactionRepository) Create(transaction *model.Transaction) error {
	return config.DB.Create(transaction).Error
}

func (r *TransactionRepository) Save(transaction *model.Transaction) error {
	return config.DB.Save(transaction).Error
}

func (r *TransactionRepository) Delete(transaction *model.Transaction) error {
	return config.DB.Delete(transaction).Error
}

func (r *TransactionRepository) GetFinanceSummary(
	userID uint,
	startDate *time.Time,
	endDate *time.Time,
) (*FinanceSummary, error) {
	var income float64
	var expense float64

	incomeQuery := config.DB.
		Model(&model.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "income")

	expenseQuery := config.DB.
		Model(&model.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "expense")

	if startDate != nil {
		incomeQuery = incomeQuery.Where("transaction_date >= ?", *startDate)
		expenseQuery = expenseQuery.Where("transaction_date >= ?", *startDate)
	}

	if endDate != nil {
		incomeQuery = incomeQuery.Where("transaction_date <= ?", *endDate)
		expenseQuery = expenseQuery.Where("transaction_date <= ?", *endDate)
	}

	if err := incomeQuery.Select("COALESCE(SUM(amount), 0)").Scan(&income).Error; err != nil {
		return nil, err
	}

	if err := expenseQuery.Select("COALESCE(SUM(amount), 0)").Scan(&expense).Error; err != nil {
		return nil, err
	}

	return &FinanceSummary{
		Income:  income,
		Expense: expense,
		Balance: income - expense,
	}, nil
}
