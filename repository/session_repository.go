package repository

import (
	"time"

	"go-todo-api/config"
	"go-todo-api/model"
)

type SessionRepository struct{}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{}
}

func (r *SessionRepository) Create(session *model.UserSession) error {
	return config.DB.Create(session).Error
}

func (r *SessionRepository) FindActiveByHash(
	hash string,
) (*model.UserSession, error) {
	var session model.UserSession

	err := config.DB.
		Where(
			"refresh_token_hash = ? AND is_revoked = ? AND expires_at > ?",
			hash,
			false,
			time.Now(),
		).
		First(&session).
		Error

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) Revoke(id uint) error {
	return config.DB.
		Model(&model.UserSession{}).
		Where("id = ?", id).
		Update("is_revoked", true).
		Error
}

func (r *SessionRepository) RevokeAllForUser(userID uint) error {
	return config.DB.
		Model(&model.UserSession{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).
		Error
}

func (r *SessionRepository) CleanupExpired() error {
	return config.DB.
		Where("expires_at < ?", time.Now()).
		Delete(&model.UserSession{}).
		Error
}
