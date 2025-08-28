package repository

import (
	"context"

	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookRepository struct {
	Repository[entity.Book]
	Log *logrus.Logger
}

func NewBookRepository(db *gorm.DB, log *logrus.Logger) *BookRepository {
	return &BookRepository{
		Repository: Repository[entity.Book]{DB: db},
		Log:        log,
	}
}

func (r *BookRepository) FindByID(ctx context.Context, db *gorm.DB, id int) (*entity.Book, error) {
	var book entity.Book
	if err := db.WithContext(ctx).Preload("Category").First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) FindByTitle(ctx context.Context, title string) (*entity.Book, error) {
	var book entity.Book
	db := r.Repository.DB
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if err := db.WithContext(ctx).Where("title = ?", title).First(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) FindByCategoryID(ctx context.Context, categoryID int) (*entity.Book, error) {
	var book entity.Book
	db := r.Repository.DB
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if err := db.WithContext(ctx).Where("category_id = ?", categoryID).First(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) CountAllBooks() (int64, error) {
	var total int64
	if err := r.DB.Model(&entity.Book{}).Count(&total).Error; err != nil {
		r.Log.Error("failed to count books: ", err)
		return 0, err
	}
	return total, nil
}

func (r *BookRepository) GetPriceStats(ctx context.Context, db *gorm.DB) (*model.BookPriceStatsResponse, error) {
	var res model.BookPriceStatsResponse

	// Query agregasi harga buku
	err := db.WithContext(ctx).
		Model(&entity.Book{}).
		Select("MAX(price) as max_price, MIN(price) as min_price, AVG(price) as avg_price").
		Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}
