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

type TransactionHandler struct {
	service  *service.TransactionService
	validate *validator.Validate
}

func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{
		service:  service.NewTransactionService(),
		validate: validator.New(),
	}
}

// GetTransactions godoc
// @Summary List transactions
// @Tags Finance
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param type query string false "income|expense"
// @Param category_id query int false "Category ID"
// @Param start_date query string false "YYYY-MM-DD"
// @Param end_date query string false "YYYY-MM-DD"
// @Success 200 {object} response.APIResponse
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
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

	var txType *string
	if t := c.Query("type"); t != "" {
		txType = &t
	}

	var categoryID *uint
	if cat := c.Query("category_id"); cat != "" {
		parsed, err := strconv.ParseUint(cat, 10, 64)
		if err != nil {
			response.BadRequest(c, "invalid category_id", nil)
			return
		}
		val := uint(parsed)
		categoryID = &val
	}

	var startDate *string
	if sd := c.Query("start_date"); sd != "" {
		startDate = &sd
	}

	var endDate *string
	if ed := c.Query("end_date"); ed != "" {
		endDate = &ed
	}

	result, err := h.service.GetTransactions(
		userID,
		page,
		limit,
		txType,
		categoryID,
		startDate,
		endDate,
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

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req dto.CreateTransactionRequest

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

	transaction, err := h.service.CreateTransaction(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"transaction created successfully",
		transaction,
	)
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	transactionID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid transaction id", nil)
		return
	}

	var req dto.UpdateTransactionRequest

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

	transaction, err := h.service.UpdateTransaction(
		userID,
		uint(transactionID),
		req,
	)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"transaction updated successfully",
		transaction,
	)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	transactionID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid transaction id", nil)
		return
	}

	err = h.service.DeleteTransaction(
		userID,
		uint(transactionID),
	)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"transaction deleted successfully",
		nil,
	)
}

func (h *TransactionHandler) GetFinanceSummary(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var startDate *string
	if sd := c.Query("start_date"); sd != "" {
		startDate = &sd
	}

	var endDate *string
	if ed := c.Query("end_date"); ed != "" {
		endDate = &ed
	}

	summary, err := h.service.GetFinanceSummary(
		userID,
		startDate,
		endDate,
	)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"finance summary fetched successfully",
		summary,
	)
}
