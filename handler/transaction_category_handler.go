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

type TransactionCategoryHandler struct {
	service  *service.TransactionCategoryService
	validate *validator.Validate
}

func NewTransactionCategoryHandler() *TransactionCategoryHandler {
	return &TransactionCategoryHandler{
		service:  service.NewTransactionCategoryService(),
		validate: validator.New(),
	}
}

func (h *TransactionCategoryHandler) GetCategories(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var categoryType *string
	if t := c.Query("type"); t != "" {
		categoryType = &t
	}

	categories, err := h.service.GetCategories(
		userID,
		categoryType,
	)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"categories fetched successfully",
		categories,
	)
}

func (h *TransactionCategoryHandler) CreateCategory(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req dto.CreateTransactionCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(
			c,
			"invalid request payload",
			err.Error(),
		)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(
			c,
			"validation failed",
			utils.FormatValidationErrors(err),
		)
		return
	}

	category, err := h.service.CreateCategory(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"category created successfully",
		category,
	)
}

func (h *TransactionCategoryHandler) UpdateCategory(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	categoryID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id", nil)
		return
	}

	var req dto.UpdateTransactionCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(
			c,
			"invalid request payload",
			err.Error(),
		)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(
			c,
			"validation failed",
			utils.FormatValidationErrors(err),
		)
		return
	}

	category, err := h.service.UpdateCategory(
		userID,
		uint(categoryID),
		req,
	)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"category updated successfully",
		category,
	)
}

func (h *TransactionCategoryHandler) DeleteCategory(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	categoryID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id", nil)
		return
	}

	err = h.service.DeleteCategory(
		userID,
		uint(categoryID),
	)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"category deleted successfully",
		nil,
	)
}
