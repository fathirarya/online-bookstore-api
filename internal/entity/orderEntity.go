package entity

import "time"

type Order struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement"`
	UserID     int       `gorm:"column:user_id;not null"`
	TotalPrice float64   `gorm:"column:total_price;type:decimal(10,2);not null"`
	Status     string    `gorm:"column:status;type:enum('PENDING','PAID','CANCELLED');default:'PENDING'"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdateAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Relations
	User       User        `gorm:"foreignKey:UserID;references:ID"`
	BookOrders []BookOrder `gorm:"foreignKey:OrderID;references:ID"`
}

func (Order) TableName() string {
	return "orders"
}
