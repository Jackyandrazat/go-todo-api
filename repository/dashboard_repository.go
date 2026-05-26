package repository

import (
	"go-todo-api/config"
	"go-todo-api/model"
	"time"
)

type DashboardRepository struct{}

type TaskStats struct {
	Total     int64 `json:"total"`
	Completed int64 `json:"completed"`
	Active    int64 `json:"active"`
}

type NoteStats struct {
	Total  int64 `json:"total"`
	Pinned int64 `json:"pinned"`
}

func NewDashboardRepository() *DashboardRepository {
	return &DashboardRepository{}
}

func (r *DashboardRepository) GetTaskStats(userID uint) (*TaskStats, error) {
	var total int64
	var completed int64
	var active int64

	if err := config.DB.
		Model(&model.Todo{}).
		Where("user_id = ?", userID).
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if err := config.DB.
		Model(&model.Todo{}).
		Where("user_id = ? AND done = ?", userID, true).
		Count(&completed).
		Error; err != nil {
		return nil, err
	}

	if err := config.DB.
		Model(&model.Todo{}).
		Where("user_id = ? AND done = ?", userID, false).
		Count(&active).
		Error; err != nil {
		return nil, err
	}

	return &TaskStats{
		Total:     total,
		Completed: completed,
		Active:    active,
	}, nil
}

func (r *DashboardRepository) GetNoteStats(userID uint) (*NoteStats, error) {
	var total int64
	var pinned int64

	if err := config.DB.
		Model(&model.Note{}).
		Where("user_id = ?", userID).
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if err := config.DB.
		Model(&model.Note{}).
		Where("user_id = ? AND is_pinned = ?", userID, true).
		Count(&pinned).
		Error; err != nil {
		return nil, err
	}

	return &NoteStats{
		Total:  total,
		Pinned: pinned,
	}, nil
}

func (r *DashboardRepository) GetRecentTasks(
	userID uint,
	limit int,
) ([]model.Todo, error) {
	var todos []model.Todo

	err := config.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&todos).
		Error

	return todos, err
}

func (r *DashboardRepository) GetRecentNotes(
	userID uint,
	limit int,
) ([]model.Note, error) {
	var notes []model.Note

	err := config.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&notes).
		Error

	return notes, err
}

func (r *DashboardRepository) GetFinanceSummary(
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
		incomeQuery = incomeQuery.Where(
			"transaction_date >= ?",
			*startDate,
		)
		expenseQuery = expenseQuery.Where(
			"transaction_date >= ?",
			*startDate,
		)
	}

	if endDate != nil {
		incomeQuery = incomeQuery.Where(
			"transaction_date <= ?",
			*endDate,
		)
		expenseQuery = expenseQuery.Where(
			"transaction_date <= ?",
			*endDate,
		)
	}

	if err := incomeQuery.
		Select("COALESCE(SUM(amount), 0)").
		Scan(&income).Error; err != nil {
		return nil, err
	}

	if err := expenseQuery.
		Select("COALESCE(SUM(amount), 0)").
		Scan(&expense).Error; err != nil {
		return nil, err
	}

	return &FinanceSummary{
		Income:  income,
		Expense: expense,
		Balance: income - expense,
	}, nil
}
