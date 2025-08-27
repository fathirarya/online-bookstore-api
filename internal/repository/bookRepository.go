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

func (r *BookRepository) FindByAuthor(ctx context.Context, author string) (*entity.Book, error) {
	var book entity.Book
	db := r.Repository.DB
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if err := db.WithContext(ctx).Where("author = ?", author).First(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) FindByPrice(ctx context.Context, price float64) (*entity.Book, error) {
	var book entity.Book
	db := r.Repository.DB
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if err := db.WithContext(ctx).Where("price = ?", price).First(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) FindByYear(ctx context.Context, year int) (*entity.Book, error) {
	var book entity.Book
	db := r.Repository.DB
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if err := db.WithContext(ctx).Where("year = ?", year).First(&book).Error; err != nil {
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

func (r *BookRepository) FindByImage(ctx context.Context, image string) (*entity.Book, error) {
	var book entity.Book
	db := r.Repository.DB
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if err := db.WithContext(ctx).Where("image_base64 = ?", image).First(&book).Error; err != nil {
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

func (r *BookRepository) GetPriceStats() (*model.BookPriceStatsResponse, error) {
	type result struct {
		Max float64
		Min float64
		Avg float64
	}

	var res result
	if err := r.DB.Model(&entity.Book{}).
		Select("MAX(price) as max, MIN(price) as min, AVG(price) as avg").
		Scan(&res).Error; err != nil {
		r.Log.Error("failed to get book price stats: ", err)
		return nil, err
	}

	return &model.BookPriceStatsResponse{
		MaxPrice: res.Max,
		MinPrice: res.Min,
		AvgPrice: res.Avg,
	}, nil
}
