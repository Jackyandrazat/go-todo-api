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

type HabitHandler struct {
	service  *service.HabitService
	validate *validator.Validate
}

func NewHabitHandler() *HabitHandler {
	return &HabitHandler{
		service:  service.NewHabitService(),
		validate: validator.New(),
	}
}

func (h *HabitHandler) GetHabits(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)
	date := c.DefaultQuery("date", "")

	habits, err := h.service.GetHabits(userID, date)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "habits fetched successfully", habits)
}

func (h *HabitHandler) CreateHabit(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req dto.CreateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload", utils.FormatValidationErrors(err))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", utils.FormatValidationErrors(err))
		return
	}

	habit, err := h.service.CreateHabit(userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "habit created successfully", habit)
}

func (h *HabitHandler) ToggleHabit(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	habitID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid habit id", nil)
		return
	}

	var req dto.ToggleHabitRequest
	_ = c.ShouldBindJSON(&req)

	isCompleted, habit, err := h.service.ToggleHabit(userID, uint(habitID), req.Date)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "habit toggled successfully", gin.H{
		"is_completed": isCompleted,
		"habit":        habit,
	})
}

func (h *HabitHandler) DeleteHabit(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	habitID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid habit id", nil)
		return
	}

	if err := h.service.DeleteHabit(userID, uint(habitID)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "habit deleted successfully", nil)
}
