package service

import (
	"errors"
	"time"

	"go-todo-api/dto"
	"go-todo-api/model"
	"go-todo-api/repository"

	"gorm.io/gorm"
)

type HabitService struct {
	repo *repository.HabitRepository
}

func NewHabitService() *HabitService {
	return &HabitService{
		repo: repository.NewHabitRepository(),
	}
}

func (s *HabitService) GetHabits(userID uint, targetDate string) ([]dto.HabitItemResponse, error) {
	if targetDate == "" {
		targetDate = time.Now().Format("2006-01-02")
	}

	habits, err := s.repo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	completedMap, err := s.repo.GetCompletedHabitIDsForDate(userID, targetDate)
	if err != nil {
		return nil, err
	}

	var res []dto.HabitItemResponse
	for _, h := range habits {
		res = append(res, dto.HabitItemResponse{
			ID:                h.ID,
			Title:             h.Title,
			Category:          h.Category,
			Icon:              h.Icon,
			Streak:            h.Streak,
			LastCompletedDate: h.LastCompletedDate,
			CreatedAt:         h.CreatedAt.Format(time.RFC3339),
			IsCompletedToday:  completedMap[h.ID],
		})
	}

	return res, nil
}

func (s *HabitService) CreateHabit(userID uint, req dto.CreateHabitRequest) (*dto.HabitItemResponse, error) {
	category := "Produktivitas"
	if req.Category != "" {
		category = req.Category
	}

	icon := "Sparkles"
	if req.Icon != "" {
		icon = req.Icon
	}

	habit := model.Habit{
		UserID:   userID,
		Title:    req.Title,
		Category: category,
		Icon:     icon,
		Streak:   0,
	}

	if err := s.repo.Create(&habit); err != nil {
		return nil, err
	}

	return &dto.HabitItemResponse{
		ID:                habit.ID,
		Title:             habit.Title,
		Category:          habit.Category,
		Icon:              habit.Icon,
		Streak:            habit.Streak,
		LastCompletedDate: habit.LastCompletedDate,
		CreatedAt:         habit.CreatedAt.Format(time.RFC3339),
		IsCompletedToday:  false,
	}, nil
}

func (s *HabitService) UpdateHabit(userID uint, habitID uint, req dto.UpdateHabitRequest) (*model.Habit, error) {
	habit, err := s.repo.FindByIDAndUserID(habitID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("habit not found")
		}
		return nil, err
	}

	if req.Title != nil {
		habit.Title = *req.Title
	}
	if req.Category != nil {
		habit.Category = *req.Category
	}
	if req.Icon != nil {
		habit.Icon = *req.Icon
	}

	if err := s.repo.Save(habit); err != nil {
		return nil, err
	}

	return habit, nil
}

func (s *HabitService) ToggleHabit(userID uint, habitID uint, date string) (bool, *dto.HabitItemResponse, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	habit, err := s.repo.FindByIDAndUserID(habitID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, errors.New("habit not found")
		}
		return false, nil, err
	}

	log, err := s.repo.FindLog(userID, habitID, date)
	isCompleted := false

	if err == nil && log != nil {
		// Sudah diceklis -> Hapus log (uncheck)
		if err := s.repo.DeleteLog(log); err != nil {
			return false, nil, err
		}
		if habit.Streak > 0 {
			habit.Streak--
		}
		if habit.LastCompletedDate != nil && *habit.LastCompletedDate == date {
			habit.LastCompletedDate = nil
		}
		isCompleted = false
	} else {
		// Belum diceklis -> Tambah log (check)
		newLog := model.HabitLog{
			UserID:  userID,
			HabitID: habitID,
			Date:    date,
		}
		if err := s.repo.CreateLog(&newLog); err != nil {
			return false, nil, err
		}
		habit.Streak++
		habit.LastCompletedDate = &date
		isCompleted = true
	}

	if err := s.repo.Save(habit); err != nil {
		return false, nil, err
	}

	return isCompleted, &dto.HabitItemResponse{
		ID:                habit.ID,
		Title:             habit.Title,
		Category:          habit.Category,
		Icon:              habit.Icon,
		Streak:            habit.Streak,
		LastCompletedDate: habit.LastCompletedDate,
		CreatedAt:         habit.CreatedAt.Format(time.RFC3339),
		IsCompletedToday:  isCompleted,
	}, nil
}

func (s *HabitService) DeleteHabit(userID uint, habitID uint) error {
	habit, err := s.repo.FindByIDAndUserID(habitID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("habit not found")
		}
		return err
	}

	return s.repo.Delete(habit)
}
