package handler

import (
	"net/http"

	"go-todo-api/response"
	"go-todo-api/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service *service.DashboardService
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		service: service.NewDashboardService(),
	}
}

func (h *DashboardHandler) GetSummary(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Error(
			c,
			http.StatusUnauthorized,
			"user context missing",
			nil,
		)
		return
	}

	userID := userIDRaw.(uint)

	var month *string
	if m := c.Query("month"); m != "" {
		month = &m
	}
	summary, err := h.service.GetSummary(
		userID,
		month,
	)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			err.Error(),
			nil,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"dashboard summary fetched successfully",
		summary,
	)
}
