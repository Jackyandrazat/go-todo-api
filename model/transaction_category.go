package model

import "gorm.io/gorm"

type TransactionCategory struct {
	gorm.Model

	UserID uint `gorm:"not null;uniqueIndex:idx_user_category_type"`

	Name string `gorm:"size:100;not null;uniqueIndex:idx_user_category_type"`

	Type string `gorm:"size:20;not null;index;uniqueIndex:idx_user_category_type"`
}
