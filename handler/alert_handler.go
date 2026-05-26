package handler

import (
	"net/http"
	"strconv"

	"go-todo-api/response"
	"go-todo-api/service"

	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	service *service.AlertService
}

func NewAlertHandler() *AlertHandler {
	return &AlertHandler{
		service: service.NewAlertService(),
	}
}

func (h *AlertHandler) GetAlerts(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	result, err := h.service.GetAll(userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"alerts fetched successfully",
		result,
	)
}

func (h *AlertHandler) DeleteAlert(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	alertIDParam := c.Param("id")
	alertID, err := strconv.ParseUint(alertIDParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid alert id", nil)
		return
	}
	err = h.service.Delete(userID, uint(alertID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"alert deleted successfully",
		nil,
	)
}

func (h *AlertHandler) MarkAsRead(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	alertIDParam := c.Param("id")
	alertID, err := strconv.ParseUint(alertIDParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid alert id", nil)
		return
	}
	_, err = h.service.MarkAsRead(userID, uint(alertID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(
		c,
		http.StatusOK,
		"alert marked as read successfully",
		nil,
	)
}
