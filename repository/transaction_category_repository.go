package repository

import (
	"go-todo-api/config"
	"go-todo-api/model"
)

type TransactionCategoryRepository struct{}

func NewTransactionCategoryRepository() *TransactionCategoryRepository {
	return &TransactionCategoryRepository{}
}

func (r *TransactionCategoryRepository) FindAllByUserID(
	userID uint,
	categoryType *string,
) ([]model.TransactionCategory, error) {
	var categories []model.TransactionCategory

	query := config.DB.
		Where("user_id = ?", userID)

	if categoryType != nil {
		query = query.Where("type = ?", *categoryType)
	}

	err := query.
		Order("name ASC").
		Find(&categories).
		Error

	return categories, err
}

func (r *TransactionCategoryRepository) FindByIDAndUserID(
	id uint,
	userID uint,
) (*model.TransactionCategory, error) {
	var category model.TransactionCategory

	err := config.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&category).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *TransactionCategoryRepository) FindByNameAndType(
	userID uint,
	name string,
	categoryType string,
) (*model.TransactionCategory, error) {
	var category model.TransactionCategory

	err := config.DB.
		Where(
			"user_id = ? AND name = ? AND type = ?",
			userID,
			name,
			categoryType,
		).
		First(&category).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *TransactionCategoryRepository) Create(
	category *model.TransactionCategory,
) error {
	return config.DB.Create(category).Error
}

func (r *TransactionCategoryRepository) Save(
	category *model.TransactionCategory,
) error {
	return config.DB.Save(category).Error
}

func (r *TransactionCategoryRepository) Delete(
	category *model.TransactionCategory,
) error {
	return config.DB.Delete(category).Error
}

func (r *TransactionCategoryRepository) FindByIDTypeAndUserID(
	id uint,
	categoryType string,
	userID uint,
) (*model.TransactionCategory, error) {
	var category model.TransactionCategory

	err := config.DB.
		Where(
			"id = ? AND type = ? AND user_id = ?",
			id,
			categoryType,
			userID,
		).
		First(&category).
		Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}
