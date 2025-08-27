package model

import "time"

type CreateOrderRequest struct {
	Items []OrderItemInput `json:"items" validate:"required,dive"`
}

// OrderItemInput adalah item buku yang diorder
type OrderItemInput struct {
	BookID   int `json:"book_id" validate:"required"`
	Quantity int `json:"quantity" validate:"required,min=1,max=5"` // max per item
}

// Response
type OrderResponse struct {
	ID         int                 `json:"id"`
	UserID     int                 `json:"user_id"`
	TotalPrice float64             `json:"total_price"`
	Status     string              `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
	Items      []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	BookID   int     `json:"book_id"`
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	SubTotal float64 `json:"sub_total"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=PENDING PAID CANCELLED"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
}
