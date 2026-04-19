package service

import (
	"go-todo-api/model"
	"go-todo-api/repository"
)

func GetTodos() []model.Todo {
	return repository.GetAll()
}

func CreateTodo(todo model.Todo) model.Todo {
	return repository.Create(todo)
}

func UpdateTodo(id uint, todo model.Todo) (model.Todo, bool) {
	return repository.Update(id, todo)
}

func DeleteTodo(id uint) bool {
	return repository.Delete(id)
}
