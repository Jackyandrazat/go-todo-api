package dto

type CreateHabitRequest struct {
	Title    string `json:"title" validate:"required,min=2,max=255"`
	Category string `json:"category" validate:"omitempty,max=100"`
	Icon     string `json:"icon" validate:"omitempty,max=50"`
}

type UpdateHabitRequest struct {
	Title    *string `json:"title" validate:"omitempty,min=2,max=255"`
	Category *string `json:"category" validate:"omitempty,max=100"`
	Icon     *string `json:"icon" validate:"omitempty,max=50"`
}

type ToggleHabitRequest struct {
	Date string `json:"date" validate:"omitempty,len=10"` // YYYY-MM-DD
}

type HabitItemResponse struct {
	ID                uint    `json:"id"`
	Title             string  `json:"title"`
	Category          string  `json:"category"`
	Icon              string  `json:"icon"`
	Streak            int     `json:"streak"`
	LastCompletedDate *string `json:"last_completed_date,omitempty"`
	CreatedAt         string  `json:"created_at"`
	IsCompletedToday  bool    `json:"is_completed_today"`
}
