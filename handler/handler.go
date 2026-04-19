package handler

import (
	"go-todo-api/model"
	"go-todo-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	c.JSON(http.StatusOK, service.GetTodos())
}

func CreateTodo(c *gin.Context) {
	var todo model.Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service.CreateTodo(todo))
}

func UpdateTodo(c *gin.Context) {
	idParam := c.Param("id")

	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	id := uint(idInt)

	var updatedTodo model.Todo
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, found := service.UpdateTodo(id, updatedTodo)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func DeleteTodo(c *gin.Context) {
	idParam := c.Param("id")

	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	id := uint(idInt)

	deleted := service.DeleteTodo(id)
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.Status(http.StatusNoContent)
	c.JSON(http.StatusOK, gin.H{
		"message": "Todo deleted successfully",
	})
}
