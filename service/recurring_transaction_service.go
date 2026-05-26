package service

import (
	"errors"
	"time"

	"go-todo-api/dto"
	"go-todo-api/model"
	"go-todo-api/repository"

	"gorm.io/gorm"
)

type RecurringTransactionService struct {
	repo         *repository.RecurringTransactionRepository
	categoryRepo *repository.TransactionCategoryRepository
}

func NewRecurringTransactionService() *RecurringTransactionService {
	return &RecurringTransactionService{
		repo:         repository.NewRecurringTransactionRepository(),
		categoryRepo: repository.NewTransactionCategoryRepository(),
	}
}

func (s *RecurringTransactionService) GetAll(
	userID uint,
) ([]model.RecurringTransaction, error) {
	return s.repo.FindAllByUserID(userID)
}

func (s *RecurringTransactionService) Create(
	userID uint,
	req dto.CreateRecurringTransactionRequest,
) (*model.RecurringTransaction, error) {
	nextRun, err := time.Parse("2006-01-02", req.NextRunAt)
	if err != nil {
		return nil, errors.New("invalid next_run_at format, use YYYY-MM-DD")
	}

	_, err = s.categoryRepo.FindByIDTypeAndUserID(
		req.CategoryID,
		req.Type,
		userID,
	)
	if err != nil {
		return nil, errors.New("invalid category for recurring transaction")
	}

	item := model.RecurringTransaction{
		UserID:     userID,
		Title:      req.Title,
		Amount:     req.Amount,
		Type:       req.Type,
		CategoryID: req.CategoryID,
		Notes:      req.Notes,
		Frequency:  req.Frequency,
		NextRunAt:  nextRun,
		IsActive:   true,
	}

	err = s.repo.Create(&item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *RecurringTransactionService) Update(
	userID uint,
	recurringID uint,
	req dto.UpdateRecurringTransactionRequest,
) (*model.RecurringTransaction, error) {
	item, err := s.repo.FindByIDAndUserID(
		recurringID,
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recurring transaction not found")
		}
		return nil, err
	}

	finalType := item.Type
	finalCategoryID := item.CategoryID

	if req.Type != nil {
		finalType = *req.Type
	}

	if req.CategoryID != nil {
		finalCategoryID = *req.CategoryID
	}

	_, err = s.categoryRepo.FindByIDTypeAndUserID(
		finalCategoryID,
		finalType,
		userID,
	)
	if err != nil {
		return nil, errors.New("invalid category for recurring transaction")
	}

	if req.Title != nil {
		item.Title = *req.Title
	}

	if req.Amount != nil {
		item.Amount = *req.Amount
	}

	if req.Type != nil {
		item.Type = *req.Type
	}

	if req.CategoryID != nil {
		item.CategoryID = *req.CategoryID
	}

	if req.Notes != nil {
		item.Notes = *req.Notes
	}

	if req.Frequency != nil {
		item.Frequency = *req.Frequency
	}

	if req.NextRunAt != nil {
		parsed, err := time.Parse(
			"2006-01-02",
			*req.NextRunAt,
		)
		if err != nil {
			return nil, errors.New(
				"invalid next_run_at format, use YYYY-MM-DD",
			)
		}
		item.NextRunAt = parsed
	}

	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}

	err = s.repo.Save(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *RecurringTransactionService) Delete(
	userID uint,
	recurringID uint,
) error {
	item, err := s.repo.FindByIDAndUserID(
		recurringID,
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("recurring transaction not found")
		}
		return err
	}

	return s.repo.Delete(item)
}
