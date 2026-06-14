package handler

import (
	"net/http"
	"strconv"

	"go-todo-api/dto"
	"go-todo-api/response"
	"go-todo-api/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BudgetHandler struct {
	service  *service.BudgetService
	validate *validator.Validate
}

func NewBudgetHandler() *BudgetHandler {
	return &BudgetHandler{
		service:  service.NewBudgetService(),
		validate: validator.New(),
	}
}

func (h *BudgetHandler) GetBudgets(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	month := c.Query("month")
	if month == "" {
		response.BadRequest(c, "month query parameter is required", nil)
		return
	}

	result, err := h.service.GetBudgets(userID, month)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "budgets fetched successfully", result)
}

func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	var req dto.CreateBudgetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	result, err := h.service.CreateBudget(userID, req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"budget created successfully",
		result,
	)
}

func (h *BudgetHandler) UpdateBudget(c *gin.Context) {
	var req dto.UpdateBudgetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	budgetID64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid budget id", nil)
		return
	}

	result, err := h.service.UpdateBudget(
		userID,
		uint(budgetID64),
		req,
	)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"budget updated successfully",
		result,
	)
}

func (h *BudgetHandler) DeleteBudget(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	budgetID64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid budget id", nil)
		return
	}

	if err := h.service.DeleteBudget(
		userID,
		uint(budgetID64),
	); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"budget deleted successfully",
		nil,
	)
}
