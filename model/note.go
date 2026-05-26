package model

import "gorm.io/gorm"

type Note struct {
	gorm.Model

	UserID   uint   `gorm:"index;not null"`
	Title    string `gorm:"size:255;not null"`
	Content  string `gorm:"type:text"`
	Category string `gorm:"size:100"`
	IsPinned bool   `gorm:"default:false"`
}
