package entity

import "time"

type User struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;size:100;not null"`
	Email     string    `gorm:"column:email;size:100;unique;not null"`
	Password  string    `gorm:"column:password;size:255;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

	// Relations
	Orders []Order `gorm:"foreignKey:UserID;references:ID"`
}

func (u *User) TableName() string {
	return "users"
}
