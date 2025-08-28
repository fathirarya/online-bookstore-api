package usecase

import (
	"context"
	"fmt"

	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/enum"
	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/model/converter"
	"github.com/fathirarya/online-bookstore-api/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	OrderRepository *repository.OrderRepository
	BookRepository  *repository.BookRepository
}

func NewOrderUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	orderRepository *repository.OrderRepository, bookRepository *repository.BookRepository) *OrderUseCase {
	return &OrderUseCase{
		DB:              db,
		Log:             logger,
		Validate:        validate,
		OrderRepository: orderRepository,
		BookRepository:  bookRepository,
	}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, req *model.CreateOrderRequest, userID int) (*model.OrderResponse, error) {
	// 1. Validasi request
	if err := uc.Validate.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "validation failed, please check your input")
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	totalQuantity := 0
	for _, item := range req.Items {
		totalQuantity += item.Quantity
	}
	if totalQuantity > 5 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "maximum 5 books per transaction")
	}

	var bookOrders []entity.BookOrder
	totalPrice := 0.0

	for _, item := range req.Items {
		var book entity.Book
		if err := uc.BookRepository.FindById(tx, &book, item.BookID); err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("book not found: %d", item.BookID))
			}
			return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}

		bookOrders = append(bookOrders, entity.BookOrder{
			BookID:   book.ID,
			Quantity: item.Quantity,
			Book:     book,
		})

		totalPrice += float64(item.Quantity) * book.Price
	}

	order := &entity.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     enum.Pending,
		BookOrders: bookOrders,
	}

	if err := uc.OrderRepository.Create(tx, order); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to create order: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create order")
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create order")
	}

	fullOrder, err := uc.OrderRepository.FindByID(uc.DB, order.ID)
	if err != nil {
		uc.Log.Error("failed to fetch full order: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch order")
	}

	return converter.OrderToResponse(fullOrder), nil
}

func (uc *OrderUseCase) PayOrder(ctx context.Context, orderID int, userID int) (*model.OrderResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	order, err := uc.OrderRepository.FindByID(tx, orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("order not found: %d", orderID))
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	if order.UserID != userID {
		return nil, fiber.NewError(fiber.StatusForbidden, "you are not allowed to pay this order")
	}

	switch order.Status {
	case enum.Cancelled:
		return nil, fiber.NewError(fiber.StatusBadRequest, "order has been cancelled")
	case enum.Paid:
		return nil, fiber.NewError(fiber.StatusBadRequest, "order has been paid")
	}

	if err := uc.OrderRepository.UpdateStatus(tx, orderID, enum.Paid); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to update order status: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update order status")
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	fullOrder, err := uc.OrderRepository.FindByID(uc.DB, orderID)
	if err != nil {
		uc.Log.Error("failed to fetch full order: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch order")
	}

	return converter.OrderToResponse(fullOrder), nil
}

func (uc *OrderUseCase) GetOrdersByUser(ctx context.Context, userID int) (*model.OrderListResponse, error) {
	tx := uc.DB.WithContext(ctx)

	orders, err := uc.OrderRepository.FindByUserID(tx, userID)
	if err != nil {
		uc.Log.Error("failed to fetch orders: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch orders")
	}

	var orderResponses []model.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, *converter.OrderToResponse(&order))
	}

	return &model.OrderListResponse{
		Orders: orderResponses,
	}, nil
}
