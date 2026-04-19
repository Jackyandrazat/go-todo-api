package repository

import (
	"go-todo-api/config"
	"go-todo-api/model"
)

var todos []model.Todo
var currenID int = 1

func GetAll() []model.Todo {
	var todos []model.Todo
	config.DB.Find(&todos)
	return todos
}

func Create(todo model.Todo) model.Todo {
	config.DB.Create(&todo)
	return todo
}

func Update(id uint, updated model.Todo) (model.Todo, bool) {
	var todo model.Todo

	result := config.DB.First(&todo, id)
	if result.Error != nil {
		return model.Todo{}, false
	}

	todo.Title = updated.Title
	todo.Done = updated.Done

	config.DB.Save(&todo)

	return todo, true
}

func Delete(id uint) bool {
	result := config.DB.Delete(&model.Todo{}, id)
	return result.RowsAffected > 0
}
