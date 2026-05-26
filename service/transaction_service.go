package service

import (
	"errors"
	"time"

	"go-todo-api/dto"
	"go-todo-api/model"
	"go-todo-api/repository"

	"gorm.io/gorm"
)

type TransactionService struct {
	repo         *repository.TransactionRepository
	categoryRepo *repository.TransactionCategoryRepository
}

type PaginatedTransactionResponse struct {
	Items      []model.Transaction `json:"items"`
	Pagination Pagination          `json:"pagination"`
}

func NewTransactionService() *TransactionService {
	return &TransactionService{
		repo:         repository.NewTransactionRepository(),
		categoryRepo: repository.NewTransactionCategoryRepository(),
	}
}

func (s *TransactionService) GetTransactions(
	userID uint,
	page int,
	limit int,
	txType *string,
	categoryID *uint,
	startDateStr *string,
	endDateStr *string,
) (*PaginatedTransactionResponse, error) {
	var startDate *time.Time
	var endDate *time.Time

	if startDateStr != nil {
		parsed, err := time.Parse("2006-01-02", *startDateStr)
		if err != nil {
			return nil, errors.New("invalid start_date format, use YYYY-MM-DD")
		}
		startDate = &parsed
	}

	if endDateStr != nil {
		parsed, err := time.Parse("2006-01-02", *endDateStr)
		if err != nil {
			return nil, errors.New("invalid end_date format, use YYYY-MM-DD")
		}
		endDate = &parsed
	}

	transactions, total, err := s.repo.FindAllByUserID(
		userID,
		page,
		limit,
		txType,
		categoryID,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &PaginatedTransactionResponse{
		Items: transactions,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *TransactionService) CreateTransaction(
	userID uint,
	req dto.CreateTransactionRequest,
) (*model.Transaction, error) {
	transactionDate, err := time.Parse(
		"2006-01-02",
		req.TransactionDate,
	)
	if err != nil {
		return nil, errors.New(
			"invalid transaction_date format, use YYYY-MM-DD",
		)
	}

	transaction := model.Transaction{
		UserID:          userID,
		Title:           req.Title,
		Amount:          req.Amount,
		Type:            req.Type,
		CategoryID:      req.CategoryID,
		Notes:           req.Notes,
		TransactionDate: transactionDate,
	}

	_, err = s.categoryRepo.FindByIDTypeAndUserID(
		req.CategoryID,
		req.Type,
		userID,
	)

	if err != nil {
		return nil, errors.New("invalid category for transaction type")
	}

	err = s.repo.Create(&transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *TransactionService) UpdateTransaction(
	userID uint,
	transactionID uint,
	req dto.UpdateTransactionRequest,
) (*model.Transaction, error) {
	transaction, err := s.repo.FindByIDAndUserID(
		transactionID,
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	finalType := transaction.Type
	finalCategoryID := transaction.CategoryID

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
		return nil, errors.New("invalid category for transaction type")
	}

	if req.Title != nil {
		transaction.Title = *req.Title
	}

	if req.Amount != nil {
		transaction.Amount = *req.Amount
	}

	if req.Type != nil {
		transaction.Type = *req.Type
	}

	if req.CategoryID != nil {
		transaction.CategoryID = *req.CategoryID
	}

	if req.Notes != nil {
		transaction.Notes = *req.Notes
	}

	if req.TransactionDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.TransactionDate)
		if err != nil {
			return nil, errors.New("invalid transaction_date format, use YYYY-MM-DD")
		}
		transaction.TransactionDate = parsed
	}

	err = s.repo.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) DeleteTransaction(
	userID uint,
	transactionID uint,
) error {
	transaction, err := s.repo.FindByIDAndUserID(
		transactionID,
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("transaction not found")
		}
		return err
	}

	return s.repo.Delete(transaction)
}

func (s *TransactionService) GetFinanceSummary(
	userID uint,
	startDateStr *string,
	endDateStr *string,
) (*repository.FinanceSummary, error) {
	var startDate *time.Time
	var endDate *time.Time

	if startDateStr != nil {
		parsed, err := time.Parse("2006-01-02", *startDateStr)
		if err != nil {
			return nil, errors.New("invalid start_date format, use YYYY-MM-DD")
		}
		startDate = &parsed
	}

	if endDateStr != nil {
		parsed, err := time.Parse("2006-01-02", *endDateStr)
		if err != nil {
			return nil, errors.New("invalid end_date format, use YYYY-MM-DD")
		}
		endDate = &parsed
	}

	return s.repo.GetFinanceSummary(
		userID,
		startDate,
		endDate,
	)
}
