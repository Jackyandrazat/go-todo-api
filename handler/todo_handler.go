package handler

import (
	"net/http"
	"strconv"

	"go-todo-api/dto"
	"go-todo-api/response"
	"go-todo-api/service"
	"go-todo-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TodoHandler struct {
	service  *service.TodoService
	validate *validator.Validate
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		service:  service.NewTodoService(),
		validate: validator.New(),
	}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	var done *bool
	if doneQuery := c.Query("done"); doneQuery != "" {
		parsed := doneQuery == "true"
		done = &parsed
	}

	todos, err := h.service.GetTodos(
		userID,
		page,
		limit,
		done,
	)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"todos fetched successfully",
		todos,
	)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req dto.CreateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request payload",
			utils.FormatValidationErrors(err),
		)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"validation failed",
			utils.FormatValidationErrors(err),
		)
		return
	}

	todo, err := h.service.CreateTodo(userID, req)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			err.Error(),
			nil,
		)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"todo created successfully",
		todo,
	)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	todoID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid todo id",
			nil,
		)
		return
	}

	var req dto.UpdateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request payload",
			utils.FormatValidationErrors(err),
		)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"validation failed",
			utils.FormatValidationErrors(err),
		)
		return
	}

	todo, err := h.service.UpdateTodo(
		userID,
		uint(todoID),
		req,
	)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			err.Error(),
			nil,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"todo updated successfully",
		todo,
	)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	todoID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid todo id",
			nil,
		)
		return
	}

	err = h.service.DeleteTodo(userID, uint(todoID))
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			err.Error(),
			nil,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"todo deleted successfully",
		nil,
	)
}
