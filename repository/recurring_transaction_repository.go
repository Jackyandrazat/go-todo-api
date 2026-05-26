package repository

import (
	"go-todo-api/config"
	"go-todo-api/model"
	"time"
)

type RecurringTransactionRepository struct{}

func NewRecurringTransactionRepository() *RecurringTransactionRepository {
	return &RecurringTransactionRepository{}
}

func (r *RecurringTransactionRepository) FindAllByUserID(
	userID uint,
) ([]model.RecurringTransaction, error) {
	var items []model.RecurringTransaction

	err := config.DB.
		Preload("Category").
		Where("user_id = ?", userID).
		Order("next_run_at ASC").
		Find(&items).
		Error

	return items, err
}

func (r *RecurringTransactionRepository) FindByIDAndUserID(
	id uint,
	userID uint,
) (*model.RecurringTransaction, error) {
	var item model.RecurringTransaction

	err := config.DB.
		Preload("Category").
		Where("id = ? AND user_id = ?", id, userID).
		First(&item).
		Error

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *RecurringTransactionRepository) Create(item *model.RecurringTransaction) error {
	return config.DB.Create(item).Error
}

func (r *RecurringTransactionRepository) Save(item *model.RecurringTransaction) error {
	return config.DB.Save(item).Error
}

func (r *RecurringTransactionRepository) Delete(item *model.RecurringTransaction) error {
	return config.DB.Delete(item).Error
}

func (r *RecurringTransactionRepository) FindDueRecurringTransactions(
	now time.Time,
) ([]model.RecurringTransaction, error) {
	var items []model.RecurringTransaction

	err := config.DB.
		Preload("Category").
		Where(
			"is_active = ? AND next_run_at <= ?",
			true,
			now,
		).
		Find(&items).
		Error

	return items, err
}

func (r *RecurringTransactionRepository) TransactionExistsForDate(
	userID uint,
	title string,
	date time.Time,
) (bool, error) {
	var count int64

	err := config.DB.
		Model(&model.Transaction{}).
		Where(
			"user_id = ? AND title = ? AND DATE(transaction_date) = DATE(?)",
			userID,
			title,
			date,
		).
		Count(&count).
		Error

	return count > 0, err
}

func (r *RecurringTransactionRepository) AdvanceNextRun(
	item *model.RecurringTransaction,
) error {
	switch item.Frequency {
	case "daily":
		item.NextRunAt = item.NextRunAt.AddDate(0, 0, 1)

	case "weekly":
		item.NextRunAt = item.NextRunAt.AddDate(0, 0, 7)

	case "monthly":
		item.NextRunAt = item.NextRunAt.AddDate(0, 1, 0)

	case "yearly":
		item.NextRunAt = item.NextRunAt.AddDate(1, 0, 0)
	}

	return config.DB.Save(item).Error
}
