package repository

import (
	"context"

	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	Repository[entity.Category]
	Log *logrus.Logger
}

func NewCategoryRepository(db *gorm.DB, log *logrus.Logger) *CategoryRepository {
	return &CategoryRepository{
		Repository: Repository[entity.Category]{DB: db},
		Log:        log,
	}
}

func (r *CategoryRepository) FindByName(ctx context.Context, name string) (*entity.Category, error) {
	var category entity.Category
	db := r.Repository.DB
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if err := db.WithContext(ctx).Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
