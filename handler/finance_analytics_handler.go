package handler

import (
	"net/http"

	"go-todo-api/response"
	"go-todo-api/service"

	"github.com/gin-gonic/gin"
)

type FinanceAnalyticsHandler struct {
	service *service.FinanceAnalyticsService
}

func NewFinanceAnalyticsHandler() *FinanceAnalyticsHandler {
	return &FinanceAnalyticsHandler{
		service: service.NewFinanceAnalyticsService(),
	}
}

// @Summary Finance analytics
// @Tags Finance
// @Security BearerAuth
// @Param month query string true "YYYY-MM"
// @Success 200 {object} response.APIResponse
// @Router /finance/analytics [get]
func (h *FinanceAnalyticsHandler) GetAnalytics(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	month := c.Query("month")
	if month == "" {
		response.BadRequest(
			c,
			"month query parameter is required",
			nil,
		)
		return
	}

	result, err := h.service.GetAnalytics(userID, month)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"finance analytics fetched successfully",
		result,
	)
}
