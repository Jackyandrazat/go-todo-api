package dto

type CreateTransactionCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Type string `json:"type" validate:"required,oneof=income expense"`
}

type UpdateTransactionCategoryRequest struct {
	Name *string `json:"name" validate:"omitempty,min=2,max=100"`
	Type *string `json:"type" validate:"omitempty,oneof=income expense"`
}
