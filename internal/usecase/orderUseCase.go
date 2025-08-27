package usecase

import (
	"context"
	"fmt"

	"github.com/fathirarya/online-bookstore-api/internal/entity"
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

	// 2. Cek total quantity maksimal 5 buku per transaksi
	totalQuantity := 0
	for _, item := range req.Items {
		totalQuantity += item.Quantity
	}
	if totalQuantity > 5 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "maximum 5 books per transaction")
	}

	// 3. Ambil data buku & hitung total price
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

		// Assign Book ke BookOrder agar preload nanti lengkap
		bookOrders = append(bookOrders, entity.BookOrder{
			BookID:   book.ID,
			Quantity: item.Quantity,
			Book:     book,
		})

		totalPrice += float64(item.Quantity) * book.Price
	}

	// 4. Buat order
	order := &entity.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "PENDING",
		BookOrders: bookOrders,
	}

	// 5. Simpan order menggunakan repository
	if err := uc.OrderRepository.Create(tx, order); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to create order: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create order")
	}

	// 6. Commit transaksi
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create order")
	}

	// 7. Ambil order lagi dengan preload Book agar response lengkap
	fullOrder, err := uc.OrderRepository.FindByID(uc.DB, order.ID)
	if err != nil {
		uc.Log.Error("failed to fetch full order: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch order")
	}

	// 8. Convert entity -> response
	return converter.OrderToResponse(fullOrder), nil
}

func (uc *OrderUseCase) PayOrder(ctx context.Context, orderID int, userID int) (*model.OrderResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Ambil order dengan preload BookOrders.Book
	order, err := uc.OrderRepository.FindByID(tx, orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("order not found: %d", orderID))
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	// 2. Cek order milik user
	if order.UserID != userID {
		return nil, fiber.NewError(fiber.StatusForbidden, "you are not allowed to pay this order")
	}

	// 3. Cek status order, hanya bisa bayar jika PENDING
	if order.Status != "PENDING" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "order is not pending")
	}

	// 4. Update status menjadi PAID
	if err := uc.OrderRepository.UpdateStatus(tx, orderID, "PAID"); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to update order status: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update order status")
	}

	// 5. Commit transaksi
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	// 6. Ambil kembali order dengan preload Book untuk response
	fullOrder, err := uc.OrderRepository.FindByID(uc.DB, orderID)
	if err != nil {
		uc.Log.Error("failed to fetch full order: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch order")
	}

	// 7. Convert entity -> response
	return converter.OrderToResponse(fullOrder), nil
}

func (uc *OrderUseCase) GetOrdersByUser(ctx context.Context, userID int) (*model.OrderListResponse, error) {
	tx := uc.DB.WithContext(ctx)

	// 1️⃣ Ambil semua order milik user
	orders, err := uc.OrderRepository.FindByUserID(tx, userID)
	if err != nil {
		uc.Log.Error("failed to fetch orders: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch orders")
	}

	// 2️⃣ Convert semua order -> OrderResponse
	var orderResponses []model.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, *converter.OrderToResponse(&order))
	}

	// 3️⃣ Bungkus dalam OrderListResponse
	return &model.OrderListResponse{
		Orders: orderResponses,
	}, nil
}
