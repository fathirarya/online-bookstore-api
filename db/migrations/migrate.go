package migrations

import (
	"log"

	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.User{},
	// &entity.Category{},
	// &entity.Book{},
	// &entity.Order{},
	// &entity.BookOrder{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("✅ Database migrated successfully")
}
