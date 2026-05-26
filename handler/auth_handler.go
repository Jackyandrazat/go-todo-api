package handler

import (
	"net/http"

	"go-todo-api/dto"
	"go-todo-api/response"
	"go-todo-api/service"
	"go-todo-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service  *service.AuthService
	validate *validator.Validate
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		service:  service.NewAuthService(),
		validate: validator.New(),
	}
}

// Register godoc
// @Summary Register user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register payload"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

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

	err := h.service.Register(req)
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
		http.StatusCreated,
		"user registered successfully",
		nil,
	)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access + refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

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

	accessToken, refreshToken, user, err := h.service.Login(req, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		if utils.Logger != nil {
			utils.Logger.Error(
				"login failed",
				zap.String("email", req.Email),
				zap.Error(err),
			)
		}

		response.Error(
			c,
			http.StatusUnauthorized,
			err.Error(),
			nil,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"login successful",
		gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user": gin.H{
				"id":         user.ID,
				"name":       user.Name,
				"username":   user.Username,
				"email":      user.Email,
				"avatar_url": user.AvatarURL,
			},
		},
	)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest

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

	accessToken, refreshToken, err := h.service.Refresh(req, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		if utils.Logger != nil {
			utils.Logger.Error(
				"token refresh failed",
				zap.String("refresh_token", req.RefreshToken),
				zap.Error(err),
			)
		}

		response.Error(
			c,
			http.StatusUnauthorized,
			err.Error(),
			nil,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"token refreshed successfully",
		gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest

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

	err := h.service.Logout(req.RefreshToken, c.Request.UserAgent(), c.ClientIP())
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
		"logout successful",
		nil,
	)
}

// Me godoc
// @Summary Current user profile
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
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

	user, err := h.service.Me(userID)
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
		"user fetched successfully",
		gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
		},
	)
}
