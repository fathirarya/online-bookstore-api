package entity

import "time"

type Category struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;size:100;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdateAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Relations
	Books []Book `gorm:"foreignKey:CategoryID;references:ID"`
}

func (Category) TableName() string {
	return "categories"
}
