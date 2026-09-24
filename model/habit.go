package model

import (
	"gorm.io/gorm"
)

type Habit struct {
	gorm.Model

	UserID            uint    `gorm:"index;not null" json:"user_id"`
	Title             string  `gorm:"size:255;not null" json:"title"`
	Category          string  `gorm:"size:100;default:'Produktivitas'" json:"category"`
	Icon              string  `gorm:"size:50;default:'Sparkles'" json:"icon"`
	Streak            int     `gorm:"default:0" json:"streak"`
	LastCompletedDate *string `gorm:"size:10" json:"last_completed_date,omitempty"`

	Logs []HabitLog `gorm:"constraint:OnDelete:CASCADE;" json:"logs,omitempty"`
}

type HabitLog struct {
	gorm.Model

	UserID  uint   `gorm:"index;not null" json:"user_id"`
	HabitID uint   `gorm:"index;not null" json:"habit_id"`
	Date    string `gorm:"size:10;index;not null" json:"date"` // Format YYYY-MM-DD
}
