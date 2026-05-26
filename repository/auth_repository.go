package repository

import (
	"go-todo-api/config"
	"go-todo-api/model"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

func (r *AuthRepository) FindByEmail(email string) (*model.User, error) {
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

func (r *AuthRepository) FindByUsername(username string) (*model.User, error) {
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

func (r *AuthRepository) FindByID(id uint) (*model.User, error) {
	var user model.User

	err := config.DB.
		First(&user, id).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) UpdateRefreshTokenHash(
	userID uint,
	refreshTokenHash *string,
) error {
	return config.DB.
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("refresh_token_hash", refreshTokenHash).
		Error
}
