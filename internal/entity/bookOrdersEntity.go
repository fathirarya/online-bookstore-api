package entity

import "time"

type BookOrder struct {
	BookID    int       `gorm:"column:book_id;primaryKey"`
	OrderID   int       `gorm:"column:order_id;primaryKey"`
	Quantity  int       `gorm:"column:quantity;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

	// Relations
	Book  Book  `gorm:"foreignKey:BookID;references:ID"`
	Order Order `gorm:"foreignKey:OrderID;references:ID"`
}

func (BookOrder) TableName() string {
	return "book_orders"
}
