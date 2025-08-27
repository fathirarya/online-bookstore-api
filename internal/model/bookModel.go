package model

import "mime/multipart"

type CreateBookRequest struct {
	Title       string                `json:"title" validate:"required,max=255"`
	Author      string                `json:"author" validate:"required,max=100"`
	Price       float64               `json:"price" validate:"required,gt=0"`
	Year        int                   `json:"year" validate:"omitempty,numeric"`
	CategoryID  int                   `json:"category_id" validate:"required"`
	ImageBase64 *multipart.FileHeader `json:"image_base64,omitempty" validate:"omitempty,base64"`
}

type BookResponse struct {
	ID         int               `json:"id"`
	Title      string            `json:"title"`
	Author     string            `json:"author"`
	Price      float64           `json:"price"`
	Year       int               `json:"year"`
	CategoryID int               `json:"category_id"`
	Category   *CategoryResponse `json:"category,omitempty"`
	ImageURL   string            `json:"image_url,omitempty"`
}

type UpdateBookRequest struct {
	Title       string                `json:"title" validate:"required,max=255"`
	Author      string                `json:"author" validate:"required,max=100"`
	Price       float64               `json:"price" validate:"required,gt=0"`
	Year        int                   `json:"year" validate:"omitempty,numeric"`
	CategoryID  int                   `json:"category_id" validate:"required"`
	ImageBase64 *multipart.FileHeader `json:"image_base64,omitempty" validate:"omitempty,base64"`
}

type BookStatsResponse struct {
	TotalBooks int `json:"total_books"`
}

type BookPriceStatsResponse struct {
	MaxPrice float64 `json:"max_price"`
	MinPrice float64 `json:"min_price"`
	AvgPrice float64 `json:"avg_price"`
}
