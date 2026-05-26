package dto

type CreateTransactionRequest struct {
	Title           string  `json:"title" validate:"required,min=2,max=255"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	Type            string  `json:"type" validate:"required,oneof=income expense"`
	CategoryID      uint    `json:"category_id" validate:"required,gt=0"`
	Notes           string  `json:"notes"`
	TransactionDate string  `json:"transaction_date" validate:"required"`
}

type UpdateTransactionRequest struct {
	Title           *string  `json:"title" validate:"omitempty,min=2,max=255"`
	Amount          *float64 `json:"amount" validate:"omitempty,gt=0"`
	Type            *string  `json:"type" validate:"omitempty,oneof=income expense"`
	CategoryID      *uint    `json:"category_id" validate:"omitempty,gt=0"`
	Notes           *string  `json:"notes"`
	TransactionDate *string  `json:"transaction_date"`
}
