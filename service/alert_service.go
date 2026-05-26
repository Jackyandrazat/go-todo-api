package service

import (
	"errors"

	"go-todo-api/model"
	"go-todo-api/repository"

	"gorm.io/gorm"
)

type AlertService struct {
	repo *repository.AlertRepository
}

func NewAlertService() *AlertService {
	return &AlertService{
		repo: repository.NewAlertRepository(),
	}
}

func (s *AlertService) GetAll(userID uint) ([]model.Alert, error) {
	return s.repo.FindAllByUserID(userID)
}

func (s *AlertService) MarkAsRead(
	userID uint,
	alertID uint,
) (*model.Alert, error) {
	alert, err := s.repo.FindByIDAndUserID(alertID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alert not found")
		}
		return nil, err
	}

	alert.IsRead = true

	err = s.repo.Save(alert)
	if err != nil {
		return nil, err
	}

	return alert, nil
}

func (s *AlertService) Delete(
	userID uint,
	alertID uint,
) error {
	alert, err := s.repo.FindByIDAndUserID(alertID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("alert not found")
		}
		return err
	}

	return s.repo.Delete(alert)
}
