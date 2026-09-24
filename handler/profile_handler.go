package handler

import (
	"net/http"

	"go-todo-api/config"
	"go-todo-api/dto"
	"go-todo-api/model"
	"go-todo-api/response"
	"go-todo-api/service"
	"go-todo-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	service  *service.ProfileService
	validate *validator.Validate
}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{
		service:  service.NewProfileService(),
		validate: validator.New(),
	}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
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

	user, err := h.service.GetProfile(userID)
	if err != nil {
		response.Error(
			c,
			http.StatusNotFound,
			err.Error(),
			nil,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"profile fetched successfully",
		gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
		},
	)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
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

	var req dto.UpdateProfileRequest

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

	user, err := h.service.UpdateProfile(userID, req)
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
		"profile updated successfully",
		gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
		},
	)
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
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

	var req dto.ChangePasswordRequest

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

	err := h.service.ChangePassword(userID, req)
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
		"password changed successfully",
		nil,
	)
}

func (h *ProfileHandler) DeleteAccount(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user context missing", nil)
		return
	}
	userID := userIDRaw.(uint)

	// Hapus seluruh data user dalam satu transaksi
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("user_id = ?", userID).Delete(&model.Todo{})
		tx.Where("user_id = ?", userID).Delete(&model.Note{})
		tx.Where("user_id = ?", userID).Delete(&model.Transaction{})
		tx.Where("user_id = ?", userID).Delete(&model.TransactionCategory{})
		tx.Where("user_id = ?", userID).Delete(&model.Budget{})
		tx.Where("user_id = ?", userID).Delete(&model.RecurringTransaction{})
		tx.Where("user_id = ?", userID).Delete(&model.Alert{})
		tx.Where("user_id = ?", userID).Delete(&model.UserSession{})
		tx.Where("user_id = ?", userID).Delete(&model.Habit{})
		tx.Where("user_id = ?", userID).Delete(&model.HabitLog{})
		return tx.Delete(&model.User{}, userID).Error
	})

	if err != nil {
		response.InternalServerError(c, "failed to delete account")
		return
	}

	response.Success(c, http.StatusOK, "account deleted successfully", nil)
}
