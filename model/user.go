package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name             string `gorm:"size:100;not null"`
	Username         string `gorm:"size:50;uniqueIndex;not null"`
	Email            string `gorm:"size:100;uniqueIndex;not null"`
	PasswordHash     string `gorm:"not null"`
	RefreshTokenHash *string
	AvatarURL        *string
}
