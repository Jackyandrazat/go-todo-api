package main

import (
	"net/http"

	"go-todo-api/config"
	"go-todo-api/handler"
	"go-todo-api/middleware"
	"go-todo-api/model"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	println("START MIGRATION...")
	config.DB.AutoMigrate(&model.Todo{})
	config.DB.AutoMigrate(&model.User{})
	println("DONE MIGRATION...")

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello dari Gin 🚀",
		})
	})

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	// protected
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("/todos", handler.GetTodos)
	auth.POST("/todos", handler.CreateTodo)
	auth.PUT("/todos/:id", handler.UpdateTodo)
	auth.DELETE("/todos/:id", handler.DeleteTodo)
	r.Run(":8080")
}
