package model

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model

	UserID uint    `gorm:"index;not null"`
	Title  string  `gorm:"size:255;not null"`
	Amount float64 `gorm:"not null"`
	Type   string  `gorm:"size:20;index;not null"`

	CategoryID uint                `gorm:"index;not null"`
	Category   TransactionCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Notes           string    `gorm:"type:text"`
	TransactionDate time.Time `gorm:"index;not null"`
}
