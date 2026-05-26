package service

import (
	"errors"
	"go-todo-api/dto"
	"go-todo-api/model"
	"go-todo-api/repository"

	"gorm.io/gorm"
)

type BudgetService struct {
	repo *repository.BudgetRepository
}

func NewBudgetService() *BudgetService {
	return &BudgetService{
		repo: repository.NewBudgetRepository(),
	}
}

func (s *BudgetService) GetBudgets(
	userID uint,
	month string,
) ([]repository.BudgetWithUsage, error) {
	return s.repo.GetBudgetsWithUsage(userID, month)
}

func (s *BudgetService) CreateBudget(
	userID uint,
	req dto.CreateBudgetRequest,
) (*model.Budget, error) {
	existing, err := s.repo.FindByCategoryAndMonth(
		userID,
		req.CategoryID,
		req.Month,
	)

	if err == nil && existing != nil {
		return nil, errors.New("budget already exists for this category/month")
	}

	budget := model.Budget{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Month:      req.Month,
	}

	err = s.repo.Create(&budget)
	if err != nil {
		return nil, err
	}

	return &budget, nil
}

func (s *BudgetService) UpdateBudget(
	userID uint,
	budgetID uint,
	req dto.UpdateBudgetRequest,
) (*model.Budget, error) {
	budget, err := s.repo.FindByIDAndUserID(budgetID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("budget not found")
		}
		return nil, err
	}

	if req.Amount != nil {
		budget.Amount = *req.Amount
	}

	if req.Month != nil {
		budget.Month = *req.Month
	}

	err = s.repo.Save(budget)
	if err != nil {
		return nil, err
	}

	return budget, nil
}

func (s *BudgetService) DeleteBudget(
	userID uint,
	budgetID uint,
) error {
	budget, err := s.repo.FindByIDAndUserID(budgetID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("budget not found")
		}
		return err
	}

	return s.repo.Delete(budget)
}
