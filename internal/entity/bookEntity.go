package entity

type Book struct {
	ID          int     `gorm:"column:id;primaryKey;autoIncrement"`
	Title       string  `gorm:"column:title;size:255;not null"`
	Author      string  `gorm:"column:author;size:100;not null"`
	Price       float64 `gorm:"column:price;type:decimal(10,2);not null"`
	Year        int     `gorm:"column:year"`
	CategoryID  int     `gorm:"column:category_id;not null"`
	ImageBase64 string  `gorm:"column:image_base64;type:text"`

	// Relations
	Category   Category    `gorm:"foreignKey:CategoryID;references:ID"`
	BookOrders []BookOrder `gorm:"foreignKey:BookID;references:ID"`
}

func (Book) TableName() string {
	return "books"
}
