package repository

import (
	"time"

	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/enum"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderRepository struct {
	CommonQuery[entity.Order]
	Log *logrus.Logger
}

func NewOrderRepository(db *gorm.DB, log *logrus.Logger) *OrderRepository {
	return &OrderRepository{
		CommonQuery: CommonQuery[entity.Order]{DB: db},
		Log:         log,
	}
}

// Create simpan order beserta BookOrders
func (r *OrderRepository) Create(tx *gorm.DB, order *entity.Order) error {
	return tx.Create(order).Error
}

// FindByID preload BookOrders dan Book
func (r *OrderRepository) FindByID(tx *gorm.DB, orderID int) (*entity.Order, error) {
	var order entity.Order
	if err := tx.Preload("BookOrders.Book").First(&order, orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByUserID preload BookOrders dan Book untuk semua order user
func (r *OrderRepository) FindByUserID(tx *gorm.DB, userID int) ([]entity.Order, error) {
	var orders []entity.Order
	if err := tx.Preload("BookOrders.Book").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) UpdateStatus(tx *gorm.DB, orderID int, status string) error {
	result := tx.Model(&entity.Order{}).Where("id = ?", orderID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *OrderRepository) CancelExpiredOrders(tx *gorm.DB) error {
	cutoff := time.Now().Add(-15 * time.Minute) // waktu 15 menit lalu

	result := tx.Model(&entity.Order{}).
		Where("status = ?", enum.Pending).
		Where("created_at <= ?", cutoff).
		Update("status", enum.Cancelled)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
