package service

import (
	"errors"

	"go-todo-api/dto"
	"go-todo-api/model"
	"go-todo-api/repository"
	"go-todo-api/utils"

	"gorm.io/gorm"
)

type ProfileService struct {
	repo *repository.ProfileRepository
}

func NewProfileService() *ProfileService {
	return &ProfileService{
		repo: repository.NewProfileRepository(),
	}
}

func (s *ProfileService) GetProfile(userID uint) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (s *ProfileService) UpdateProfile(
	userID uint,
	req dto.UpdateProfileRequest,
) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.Email != nil {
		existing, err := s.repo.FindByEmail(*req.Email)
		if err == nil && existing.ID != user.ID {
			return nil, errors.New("email already registered")
		}
		user.Email = *req.Email
	}

	if req.Username != nil {
		existing, err := s.repo.FindByUsername(*req.Username)
		if err == nil && existing.ID != user.ID {
			return nil, errors.New("username already taken")
		}
		user.Username = *req.Username
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	err = s.repo.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *ProfileService) ChangePassword(
	userID uint,
	req dto.ChangePasswordRequest,
) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if !utils.CheckPasswordHash(
		req.CurrentPassword,
		user.PasswordHash,
	) {
		return errors.New("current password is incorrect")
	}

	if utils.CheckPasswordHash(
		req.NewPassword,
		user.PasswordHash,
	) {
		return errors.New("new password cannot be the same as current password")
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash

	err = s.repo.Save(user)
	if err != nil {
		return err
	}

	return nil
}
