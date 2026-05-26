package model

import (
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	gorm.Model

	UserID uint `gorm:"index;not null"`

	RefreshTokenHash string `gorm:"size:255;not null"`

	UserAgent string `gorm:"size:500"`
	IPAddress string `gorm:"size:100"`

	ExpiresAt time.Time `gorm:"index;not null"`
	IsRevoked bool      `gorm:"default:false"`
}
