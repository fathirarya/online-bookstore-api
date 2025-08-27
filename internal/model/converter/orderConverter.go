package converter

import (
	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/model"
)

func BookOrderToResponse(bookOrder *entity.BookOrder, book *entity.Book) model.OrderItemResponse {
	return model.OrderItemResponse{
		BookID:   book.ID,
		Title:    book.Title,
		Price:    book.Price,
		Quantity: bookOrder.Quantity,
		SubTotal: float64(bookOrder.Quantity) * book.Price,
	}
}

// OrderToResponse converts an Order entity to OrderResponse
func OrderToResponse(order *entity.Order) *model.OrderResponse {
	var items []model.OrderItemResponse

	for _, bo := range order.BookOrders {
		// Asumsikan Book sudah dipreload di entity.Order.BookOrders[i].Book
		items = append(items, BookOrderToResponse(&bo, &bo.Book))
	}

	return &model.OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
		Items:      items,
	}
}
