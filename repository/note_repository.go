package repository

import (
	"go-todo-api/config"
	"go-todo-api/model"
)

type NoteRepository struct{}

func NewNoteRepository() *NoteRepository {
	return &NoteRepository{}
}

func (r *NoteRepository) FindAllByUserID(
	userID uint,
	page int,
	limit int,
	pinned *bool,
) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	query := config.DB.
		Model(&model.Note{}).
		Where("user_id = ?", userID)

	if pinned != nil {
		query = query.Where("is_pinned = ?", *pinned)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := query.
		Order("is_pinned DESC").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notes).
		Error

	return notes, total, err
}

func (r *NoteRepository) FindByIDAndUserID(
	id uint,
	userID uint,
) (*model.Note, error) {
	var note model.Note

	err := config.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&note).
		Error

	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *NoteRepository) Create(note *model.Note) error {
	return config.DB.Create(note).Error
}

func (r *NoteRepository) Save(note *model.Note) error {
	return config.DB.Save(note).Error
}

func (r *NoteRepository) Delete(note *model.Note) error {
	return config.DB.Delete(note).Error
}
