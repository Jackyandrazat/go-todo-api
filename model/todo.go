package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model

	UserID      uint   `gorm:"index;not null"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	Priority    string `gorm:"size:20;default:'medium'"`
	Category    string `gorm:"size:100"`
	Done        bool   `gorm:"default:false"`
	DueDate     *time.Time
}
