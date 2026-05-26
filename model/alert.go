package model

import "gorm.io/gorm"

type Alert struct {
	gorm.Model

	UserID uint `gorm:"index;not null"`

	Type    string `gorm:"size:50;index;not null"`
	Title   string `gorm:"size:255;not null"`
	Message string `gorm:"type:text;not null"`

	IsRead bool `gorm:"default:false"`
}
