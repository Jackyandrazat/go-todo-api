package repository

import (
	"go-todo-api/config"
	"go-todo-api/model"
)

type ProfileRepository struct{}

func NewProfileRepository() *ProfileRepository {
	return &ProfileRepository{}
}

func (r *ProfileRepository) FindByID(userID uint) (*model.User, error) {
	var user model.User

	err := config.DB.
		First(&user, userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ProfileRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := config.DB.
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ProfileRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User

	err := config.DB.
		Where("username = ?", username).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ProfileRepository) Save(user *model.User) error {
	return config.DB.Save(user).Error
}
