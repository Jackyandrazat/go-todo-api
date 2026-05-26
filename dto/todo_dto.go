package dto

type CreateTodoRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	Priority    string  `json:"priority" validate:"omitempty,oneof=low medium high"`
	Category    string  `json:"category"`
	DueDate     *string `json:"due_date"`
}

type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *string `json:"priority" validate:"omitempty,oneof=low medium high"`
	Category    *string `json:"category"`
	Done        *bool   `json:"done"`
	DueDate     *string `json:"due_date"`
}
