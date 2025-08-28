package usecase

import (
	"context"

	"github.com/fathirarya/online-bookstore-api/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderCronJob struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	OrderRepository *repository.OrderRepository
}

func NewOrderCronJob(db *gorm.DB, logger *logrus.Logger,
	orderRepository *repository.OrderRepository) *OrderCronJob {
	return &OrderCronJob{
		DB:              db,
		Log:             logger,
		OrderRepository: orderRepository,
	}
}

func (w OrderCronJob) CheckingOrderPaymentStatus(ctx context.Context) error {
	w.Log.Info("cron job started")

	if err := w.OrderRepository.CancelExpiredOrders(w.DB.WithContext(ctx)); err != nil {
		w.Log.Error("failed to cancel expired orders: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to cancel expired orders")
	}

	w.Log.Info("cron job done")
	return nil
}
