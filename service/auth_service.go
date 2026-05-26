package service

import (
	"errors"
	"time"

	"go-todo-api/dto"
	"go-todo-api/model"
	"go-todo-api/repository"
	"go-todo-api/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	repo        *repository.AuthRepository
	sessionRepo *repository.SessionRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo:        repository.NewAuthRepository(),
		sessionRepo: repository.NewSessionRepository(),
	}
}

func (s *AuthService) Register(
	req dto.RegisterRequest,
) error {
	existingEmail, _ := s.repo.FindByEmail(req.Email)
	if existingEmail != nil {
		return errors.New("email already registered")
	}

	existingUsername, _ := s.repo.FindByUsername(req.Username)
	if existingUsername != nil {
		return errors.New("username already taken")
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := model.User{
		Name:         req.Name,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	return s.repo.CreateUser(&user)
}

func (s *AuthService) Login(
	req dto.LoginRequest,
	userAgent string,
	ip string,
) (string, string, *model.User, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(
		req.Password,
		user.PasswordHash,
	) {
		return "", "", nil, errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
	)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", "", nil, err
	}
	refreshHash := utils.HashToken(refreshToken)

	session := model.UserSession{
		UserID:           user.ID,
		RefreshTokenHash: refreshHash,
		UserAgent:        userAgent,
		IPAddress:        ip,
		ExpiresAt:        time.Now().Add(30 * 24 * time.Hour),
		IsRevoked:        false,
	}

	err = s.sessionRepo.Create(&session)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *AuthService) Refresh(
	req dto.RefreshRequest,
	userAgent string,
	ip string,
) (string, string, error) {
	refreshHash := utils.HashToken(req.RefreshToken)

	session, err := s.sessionRepo.FindActiveByHash(
		refreshHash,
	)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	user, err := s.repo.FindByID(session.UserID)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	err = s.sessionRepo.Revoke(session.ID)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
	)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", "", err
	}
	newRefreshHash := utils.HashToken(newRefreshToken)

	newSession := model.UserSession{
		UserID:           user.ID,
		RefreshTokenHash: newRefreshHash,
		UserAgent:        userAgent,
		IPAddress:        ip,
		ExpiresAt:        time.Now().Add(30 * 24 * time.Hour),
		IsRevoked:        false,
	}

	err = s.sessionRepo.Create(&newSession)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) Logout(
	refreshToken string,
	userAgent string,
	ip string,
) error {
	refreshHash := utils.HashToken(refreshToken)

	session, err := s.sessionRepo.FindActiveByHash(
		refreshHash,
	)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	return s.sessionRepo.Revoke(session.ID)
}

func (s *AuthService) Me(
	userID uint,
) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
