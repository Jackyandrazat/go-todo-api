package handler

import (
	"net/http"

	"go-todo-api/dto"
	"go-todo-api/response"
	"go-todo-api/service"
	"go-todo-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
