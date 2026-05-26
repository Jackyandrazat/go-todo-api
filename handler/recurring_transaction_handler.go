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

type RecurringTransactionHandler struct {
	service  *service.RecurringTransactionService
	validate *validator.Validate
}

func NewRecurringTransactionHandler() *RecurringTransactionHandler {
	return &RecurringTransactionHandler{
		service:  service.NewRecurringTransactionService(),
		validate: validator.New(),
	}
}

func (h *RecurringTransactionHandler) GetRecurringTransactions(c *gin.Context) {
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

	result, err := h.service.GetAll(
		userID,
	)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"transactions fetched successfully",
		result,
	)
}

func (h *RecurringTransactionHandler) CreateRecurringTransaction(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req dto.CreateRecurringTransactionRequest

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

	recurringTransaction, err := h.service.Create(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"recurring transaction created successfully",
		recurringTransaction,
	)
}

func (h *RecurringTransactionHandler) UpdateRecurringTransaction(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	recurringTransactionID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid transaction id", nil)
		return
	}

	var req dto.UpdateRecurringTransactionRequest

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

	recurringTransaction, err := h.service.Update(
		userID,
		uint(recurringTransactionID),
		req,
	)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"recurring transaction updated successfully",
		recurringTransaction,
	)
}

func (h *RecurringTransactionHandler) DeleteRecurringTransaction(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	recurringTransactionID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid recurring transaction id", nil)
		return
	}

	err = h.service.Delete(
		userID,
		uint(recurringTransactionID),
	)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"recurring transaction deleted successfully",
		nil,
	)
}
