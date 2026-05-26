package dto

type UpdateProfileRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=2,max=100"`
	Username *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email" validate:"omitempty,email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}
