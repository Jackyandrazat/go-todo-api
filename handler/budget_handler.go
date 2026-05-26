package handler

import (
	"net/http"

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

func (h *BudgetHandler) CreateBudget(c *gin.Context) { /* same pattern */ }
func (h *BudgetHandler) UpdateBudget(c *gin.Context) { /* same pattern */ }
func (h *BudgetHandler) DeleteBudget(c *gin.Context) { /* same pattern */ }
