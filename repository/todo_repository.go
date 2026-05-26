package repository

import (
	"go-todo-api/config"
	"go-todo-api/model"
)

type TodoRepository struct{}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{}
}

func (r *TodoRepository) FindAllByUserID(
	userID uint,
	page int,
	limit int,
	done *bool,
) ([]model.Todo, int64, error) {
	var todos []model.Todo
	var total int64

	query := config.DB.
		Model(&model.Todo{}).
		Where("user_id = ?", userID)

	if done != nil {
		query = query.Where("done = ?", *done)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&todos).
		Error

	return todos, total, err
}

func (r *TodoRepository) FindByIDAndUserID(
	id uint,
	userID uint,
) (*model.Todo, error) {
	var todo model.Todo

	err := config.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&todo).
		Error

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepository) Create(todo *model.Todo) error {
	return config.DB.Create(todo).Error
}

func (r *TodoRepository) Save(todo *model.Todo) error {
	return config.DB.Save(todo).Error
}

func (r *TodoRepository) Delete(todo *model.Todo) error {
	return config.DB.Delete(todo).Error
}
