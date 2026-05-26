package model

import "gorm.io/gorm"

type Budget struct {
	gorm.Model

	UserID uint `gorm:"not null;uniqueIndex:idx_budget_user_category_month"`

	CategoryID uint                `gorm:"not null;uniqueIndex:idx_budget_user_category_month"`
	Category   TransactionCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Amount float64 `gorm:"not null"`

	Month string `gorm:"size:7;not null;uniqueIndex:idx_budget_user_category_month"`
}
