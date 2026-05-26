package dto

type CreateBudgetRequest struct {
	CategoryID uint    `json:"category_id" validate:"required,gt=0"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	Month      string  `json:"month" validate:"required"`
}

type UpdateBudgetRequest struct {
	Amount *float64 `json:"amount" validate:"omitempty,gt=0"`
	Month  *string  `json:"month"`
}
