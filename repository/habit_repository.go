package repository

import (
	"go-todo-api/config"
	"go-todo-api/model"
)

type HabitRepository struct{}

func NewHabitRepository() *HabitRepository {
	return &HabitRepository{}
}

func (r *HabitRepository) FindAllByUserID(userID uint) ([]model.Habit, error) {
	var habits []model.Habit
	err := config.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&habits).
		Error
	return habits, err
}

func (r *HabitRepository) FindByIDAndUserID(id uint, userID uint) (*model.Habit, error) {
	var habit model.Habit
	err := config.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&habit).
		Error
	if err != nil {
		return nil, err
	}
	return &habit, nil
}

func (r *HabitRepository) Create(habit *model.Habit) error {
	return config.DB.Create(habit).Error
}

func (r *HabitRepository) Save(habit *model.Habit) error {
	return config.DB.Save(habit).Error
}

func (r *HabitRepository) Delete(habit *model.Habit) error {
	return config.DB.Delete(habit).Error
}

func (r *HabitRepository) FindLog(userID uint, habitID uint, date string) (*model.HabitLog, error) {
	var log model.HabitLog
	err := config.DB.
		Where("user_id = ? AND habit_id = ? AND date = ?", userID, habitID, date).
		First(&log).
		Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *HabitRepository) CreateLog(log *model.HabitLog) error {
	return config.DB.Create(log).Error
}

func (r *HabitRepository) DeleteLog(log *model.HabitLog) error {
	return config.DB.Delete(log).Error
}

func (r *HabitRepository) GetCompletedHabitIDsForDate(userID uint, date string) (map[uint]bool, error) {
	var logs []model.HabitLog
	err := config.DB.
		Where("user_id = ? AND date = ?", userID, date).
		Find(&logs).
		Error
	if err != nil {
		return nil, err
	}

	completed := make(map[uint]bool)
	for _, l := range logs {
		completed[l.HabitID] = true
	}
	return completed, nil
}
