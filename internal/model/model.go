package model

// WebResponse is a generic API response wrapper
type WebResponse[T any] struct {
	Page       int               `json:"page,omitempty"`
	Size       int               `json:"size,omitempty"`
	TotalItems int64             `json:"total_items,omitempty"`
	TotalPages int64             `json:"total_pages,omitempty"`
	Data       T                 `json:"data,omitempty"`
	Errors     map[string]string `json:"errors,omitempty"`
	Message    string            `json:"message,omitempty"`
}
