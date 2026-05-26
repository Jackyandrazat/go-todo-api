package repository

import (
	"go-todo-api/config"
	"go-todo-api/model"
)

type BudgetRepository struct{}

type BudgetWithUsage struct {
	ID           uint    `json:"id"`
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"budget"`
	Spent        float64 `json:"spent"`
	Remaining    float64 `json:"remaining"`
	Month        string  `json:"month"`
}

func NewBudgetRepository() *BudgetRepository {
	return &BudgetRepository{}
}

func (r *BudgetRepository) FindByIDAndUserID(
	id uint,
	userID uint,
) (*model.Budget, error) {
	var budget model.Budget

	err := config.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&budget).
		Error

	if err != nil {
		return nil, err
	}

	return &budget, nil
}

func (r *BudgetRepository) FindByCategoryAndMonth(
	userID uint,
	categoryID uint,
	month string,
) (*model.Budget, error) {
	var budget model.Budget

	err := config.DB.
		Where(
			"user_id = ? AND category_id = ? AND month = ?",
			userID,
			categoryID,
			month,
		).
		First(&budget).
		Error

	if err != nil {
		return nil, err
	}

	return &budget, nil
}

func (r *BudgetRepository) Create(budget *model.Budget) error {
	return config.DB.Create(budget).Error
}

func (r *BudgetRepository) Save(budget *model.Budget) error {
	return config.DB.Save(budget).Error
}

func (r *BudgetRepository) Delete(budget *model.Budget) error {
	return config.DB.Delete(budget).Error
}

func (r *BudgetRepository) GetBudgetsWithUsage(
	userID uint,
	month string,
) ([]BudgetWithUsage, error) {
	var result []BudgetWithUsage

	err := config.DB.Raw(`
		SELECT
			b.id,
			b.category_id,
			tc.name as category_name,
			b.amount,
			COALESCE(SUM(t.amount), 0) as spent,
			(b.amount - COALESCE(SUM(t.amount), 0)) as remaining,
			b.month
		FROM budgets b
		JOIN transaction_categories tc
			ON tc.id = b.category_id
		LEFT JOIN transactions t
			ON t.user_id = b.user_id
			AND t.category_id = tc.id
			AND t.type = 'expense'
			AND TO_CHAR(t.transaction_date, 'YYYY-MM') = b.month
		WHERE b.user_id = ?
			AND b.month = ?
		GROUP BY b.id, tc.name
		ORDER BY tc.name ASC
	`, userID, month).Scan(&result).Error

	return result, err
}
