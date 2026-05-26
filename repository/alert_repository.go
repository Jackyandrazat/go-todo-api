package repository

import (
	"go-todo-api/config"
	"go-todo-api/model"
)

type AlertRepository struct{}

func NewAlertRepository() *AlertRepository {
	return &AlertRepository{}
}

func (r *AlertRepository) FindAllByUserID(
	userID uint,
) ([]model.Alert, error) {
	var alerts []model.Alert

	err := config.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&alerts).
		Error

	return alerts, err
}

func (r *AlertRepository) FindByIDAndUserID(
	id uint,
	userID uint,
) (*model.Alert, error) {
	var alert model.Alert

	err := config.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&alert).
		Error

	if err != nil {
		return nil, err
	}

	return &alert, nil
}

func (r *AlertRepository) Create(alert *model.Alert) error {
	return config.DB.Create(alert).Error
}

func (r *AlertRepository) Save(alert *model.Alert) error {
	return config.DB.Save(alert).Error
}

func (r *AlertRepository) Delete(alert *model.Alert) error {
	return config.DB.Delete(alert).Error
}

func (r *AlertRepository) ExistsSameUnread(
	userID uint,
	alertType string,
	title string,
) (bool, error) {
	var count int64

	err := config.DB.
		Model(&model.Alert{}).
		Where(
			"user_id = ? AND type = ? AND title = ? AND is_read = ?",
			userID,
			alertType,
			title,
			false,
		).
		Count(&count).
		Error

	return count > 0, err
}
