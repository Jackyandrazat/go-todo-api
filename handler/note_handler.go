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

type NoteHandler struct {
	service  *service.NoteService
	validate *validator.Validate
}

func NewNoteHandler() *NoteHandler {
	return &NoteHandler{
		service:  service.NewNoteService(),
		validate: validator.New(),
	}
}

func (h *NoteHandler) GetNotes(c *gin.Context) {
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

	var pinned *bool
	if pinnedQuery := c.Query("pinned"); pinnedQuery != "" {
		parsed := pinnedQuery == "true"
		pinned = &parsed
	}

	notes, err := h.service.GetNotes(
		userID,
		page,
		limit,
		pinned,
	)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"notes fetched successfully",
		notes,
	)
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req dto.CreateNoteRequest

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

	note, err := h.service.CreateNote(userID, req)
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
		"note created successfully",
		note,
	)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	noteID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid note id",
			nil,
		)
		return
	}

	var req dto.UpdateNoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request payload",
			utils.FormatValidationErrors(err),
		)
		return
	}

	note, err := h.service.UpdateNote(
		userID,
		uint(noteID),
		req,
	)
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
		"note updated successfully",
		note,
	)
}

func (h *NoteHandler) TogglePin(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	noteID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid note id",
			nil,
		)
		return
	}

	var req dto.TogglePinRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request payload",
			utils.FormatValidationErrors(err),
		)
		return
	}

	note, err := h.service.TogglePin(
		userID,
		uint(noteID),
		req,
	)
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
		"note pin updated successfully",
		note,
	)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	idParam := c.Param("id")
	noteID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid note id",
			nil,
		)
		return
	}

	err = h.service.DeleteNote(userID, uint(noteID))
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
		"note deleted successfully",
		nil,
	)
}
