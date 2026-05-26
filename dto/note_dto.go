package dto

type CreateNoteRequest struct {
	Title    string `json:"title" validate:"required,min=1,max=255"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type UpdateNoteRequest struct {
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	Category *string `json:"category"`
}

type TogglePinRequest struct {
	IsPinned bool `json:"is_pinned"`
}
