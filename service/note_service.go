package service

import (
	"errors"

	"go-todo-api/dto"
	"go-todo-api/model"
	"go-todo-api/repository"

	"gorm.io/gorm"
)

type NoteService struct {
	repo *repository.NoteRepository
}

type PaginatedNoteResponse struct {
	Items      []model.Note `json:"items"`
	Pagination Pagination   `json:"pagination"`
}

func NewNoteService() *NoteService {
	return &NoteService{
		repo: repository.NewNoteRepository(),
	}
}

func (s *NoteService) GetNotes(
	userID uint,
	page int,
	limit int,
	pinned *bool,
) (*PaginatedNoteResponse, error) {
	notes, total, err := s.repo.FindAllByUserID(
		userID,
		page,
		limit,
		pinned,
	)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &PaginatedNoteResponse{
		Items: notes,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *NoteService) CreateNote(
	userID uint,
	req dto.CreateNoteRequest,
) (*model.Note, error) {
	note := model.Note{
		UserID:   userID,
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		IsPinned: false,
	}

	err := s.repo.Create(&note)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (s *NoteService) UpdateNote(
	userID uint,
	noteID uint,
	req dto.UpdateNoteRequest,
) (*model.Note, error) {
	note, err := s.repo.FindByIDAndUserID(noteID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("note not found")
		}
		return nil, err
	}

	if req.Title != nil {
		note.Title = *req.Title
	}

	if req.Content != nil {
		note.Content = *req.Content
	}

	if req.Category != nil {
		note.Category = *req.Category
	}

	err = s.repo.Save(note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteService) TogglePin(
	userID uint,
	noteID uint,
	req dto.TogglePinRequest,
) (*model.Note, error) {
	note, err := s.repo.FindByIDAndUserID(noteID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("note not found")
		}
		return nil, err
	}

	note.IsPinned = req.IsPinned

	err = s.repo.Save(note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteService) DeleteNote(
	userID uint,
	noteID uint,
) error {
	note, err := s.repo.FindByIDAndUserID(noteID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("note not found")
		}
		return err
	}

	return s.repo.Delete(note)
}
