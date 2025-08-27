package usecase

import (
	"context"

	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/model/converter"
	"github.com/fathirarya/online-bookstore-api/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	CategoryRepository *repository.CategoryRepository
}

func NewCategoryUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	categoryRepository *repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		DB:                 db,
		Log:                logger,
		Validate:           validate,
		CategoryRepository: categoryRepository,
	}
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, req *model.CreateCategoryRequest) (*model.CreateCategoryResponse, error) {
	if err := uc.Validate.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "validation failed, please check your input")
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cek name category sudah terdaftar
	existingCategory, err := uc.CategoryRepository.FindByName(ctx, req.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	if existingCategory != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "category already exists")
	}

	category := &entity.Category{
		Name: req.Name,
	}

	if err := uc.CategoryRepository.Create(tx, category); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to create category: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create category")
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create category")
	}

	return converter.CategoryToResponse(category), nil
}

func (uc *CategoryUseCase) ListCategories(ctx context.Context, page, size int) ([]*model.CategoryResponse, *model.PageMetadata, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	var categories []entity.Category
	total, err := uc.CategoryRepository.Paginate(ctx, uc.DB, page, size, &categories)
	if err != nil {
		uc.Log.Error("failed to list categories: ", err)
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, "failed to list categories")
	}

	response := make([]*model.CategoryResponse, len(categories))
	for i, c := range categories {
		response[i] = &model.CategoryResponse{
			ID:   c.ID,
			Name: c.Name,
		}
	}

	totalPage := (total + int64(size) - 1) / int64(size)
	pageMeta := &model.PageMetadata{
		Page:      page,
		Size:      size,
		TotalItem: total,
		TotalPage: totalPage,
	}

	return response, pageMeta, nil
}

func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, id int, req *model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	// Validasi request
	if err := uc.Validate.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Mulai transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Ambil data lama
	var category entity.Category
	if err := tx.First(&category, id).Error; err != nil {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusNotFound, "category not found")
	}

	// Cek nama baru sudah ada atau belum (optional)
	existingCategory, err := uc.CategoryRepository.FindByName(ctx, req.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	if existingCategory != nil && existingCategory.ID != id {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusConflict, "category name already exists")
	}

	// Update field
	category.Name = req.Name

	// Simpan perubahan
	if err := uc.CategoryRepository.Update(tx, &category); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to update category: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update category")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update category")
	}

	// Mapping ke response
	resp := &model.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	return resp, nil
}

func (uc *CategoryUseCase) DeleteCategory(ctx context.Context, id int) error {
	// Mulai transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Ambil data lama
	var category entity.Category
	if err := tx.First(&category, id).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusNotFound, "category not found")
	}

	// Hapus data
	if err := uc.CategoryRepository.Delete(tx, &category); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to delete category: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete category")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete category")
	}

	return nil
}
