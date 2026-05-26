package service

import (
	"errors"
	"time"

	"go-todo-api/dto"
	"go-todo-api/model"
	"go-todo-api/repository"

	"gorm.io/gorm"
)

type TodoService struct {
	repo *repository.TodoRepository
}

type PaginatedTodoResponse struct {
	Items      []model.Todo `json:"items"`
	Pagination Pagination   `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewTodoService() *TodoService {
	return &TodoService{
		repo: repository.NewTodoRepository(),
	}
}

func (s *TodoService) GetTodos(
	userID uint,
	page int,
	limit int,
	done *bool,
) (*PaginatedTodoResponse, error) {
	todos, total, err := s.repo.FindAllByUserID(
		userID,
		page,
		limit,
		done,
	)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &PaginatedTodoResponse{
		Items: todos,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *TodoService) CreateTodo(
	userID uint,
	req dto.CreateTodoRequest,
) (*model.Todo, error) {
	var dueDate *time.Time

	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, errors.New("invalid due_date format, use YYYY-MM-DD")
		}
		dueDate = &parsed
	}

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	todo := model.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
		Category:    req.Category,
		DueDate:     dueDate,
		Done:        false,
	}

	err := s.repo.Create(&todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *TodoService) UpdateTodo(
	userID uint,
	todoID uint,
	req dto.UpdateTodoRequest,
) (*model.Todo, error) {
	todo, err := s.repo.FindByIDAndUserID(todoID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}

	if req.Description != nil {
		todo.Description = *req.Description
	}

	if req.Priority != nil {
		todo.Priority = *req.Priority
	}

	if req.Category != nil {
		todo.Category = *req.Category
	}

	if req.Done != nil {
		todo.Done = *req.Done
	}

	if req.DueDate != nil {
		if *req.DueDate == "" {
			todo.DueDate = nil
		} else {
			parsed, err := time.Parse("2006-01-02", *req.DueDate)
			if err != nil {
				return nil, errors.New("invalid due_date format")
			}
			todo.DueDate = &parsed
		}
	}

	err = s.repo.Save(todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) DeleteTodo(
	userID uint,
	todoID uint,
) error {
	todo, err := s.repo.FindByIDAndUserID(todoID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("todo not found")
		}
		return err
	}

	return s.repo.Delete(todo)
}
