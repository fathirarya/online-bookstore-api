package entity

type Category struct {
	ID   int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name;size:100;not null"`

	// Relations
	Books []Book `gorm:"foreignKey:CategoryID;references:ID"`
}

func (Category) TableName() string {
	return "categories"
}
