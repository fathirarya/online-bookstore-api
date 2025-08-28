package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type CommonQuery[T any] struct {
	DB *gorm.DB
}

func (r *CommonQuery[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *CommonQuery[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *CommonQuery[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (r *CommonQuery[T]) FindById(db *gorm.DB, entity *T, id any) error {
	return db.Where("id = ?", id).Take(entity).Error
}

func (r *CommonQuery[T]) Paginate(ctx context.Context, db *gorm.DB, page, size int, entities *[]T) (total int64, err error) {
	// Validasi page & size
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	// Hitung total data
	if err = db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to count total: %w", err)
	}

	// Ambil data dengan limit & offset
	if err = db.WithContext(ctx).Limit(size).Offset(offset).Find(entities).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch data: %w", err)
	}

	return total, nil
}
