package service

import (
	"errors"

	"go-todo-api/dto"
	"go-todo-api/model"
	"go-todo-api/repository"

	"gorm.io/gorm"
)

type TransactionCategoryService struct {
	repo *repository.TransactionCategoryRepository
}

func NewTransactionCategoryService() *TransactionCategoryService {
	return &TransactionCategoryService{
		repo: repository.NewTransactionCategoryRepository(),
	}
}

func (s *TransactionCategoryService) GetCategories(
	userID uint,
	categoryType *string,
) ([]model.TransactionCategory, error) {
	return s.repo.FindAllByUserID(userID, categoryType)
}

func (s *TransactionCategoryService) CreateCategory(
	userID uint,
	req dto.CreateTransactionCategoryRequest,
) (*model.TransactionCategory, error) {
	existing, err := s.repo.FindByNameAndType(
		userID,
		req.Name,
		req.Type,
	)

	if err == nil && existing != nil {
		return nil, errors.New("category already exists")
	}

	category := model.TransactionCategory{
		UserID: userID,
		Name:   req.Name,
		Type:   req.Type,
	}

	err = s.repo.Create(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *TransactionCategoryService) UpdateCategory(
	userID uint,
	categoryID uint,
	req dto.UpdateTransactionCategoryRequest,
) (*model.TransactionCategory, error) {
	category, err := s.repo.FindByIDAndUserID(categoryID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	newName := category.Name
	newType := category.Type

	if req.Name != nil {
		newName = *req.Name
	}

	if req.Type != nil {
		newType = *req.Type
	}

	existing, err := s.repo.FindByNameAndType(
		userID,
		newName,
		newType,
	)

	if err == nil &&
		existing != nil &&
		existing.ID != category.ID {
		return nil, errors.New("category already exists")
	}

	category.Name = newName
	category.Type = newType

	err = s.repo.Save(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *TransactionCategoryService) DeleteCategory(
	userID uint,
	categoryID uint,
) error {
	category, err := s.repo.FindByIDAndUserID(categoryID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	err = s.repo.Delete(category)
	if err != nil {
		return errors.New(
			"cannot delete category because it is used by transactions, budgets, or recurring transactions",
		)
	}

	return nil
}
