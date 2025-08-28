package usecase

import (
	"context"

	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/model/converter"
	"github.com/fathirarya/online-bookstore-api/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	CategoryRepository *repository.CategoryRepository
}

func NewCategoryUseCase(db *gorm.DB, logger *logrus.Logger, categoryRepository *repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		DB:                 db,
		Log:                logger,
		CategoryRepository: categoryRepository,
	}
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, req *model.CreateCategoryRequest) (*model.CreateCategoryResponse, error) {
	// Start transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if category name already exists
	existingCategory, err := uc.CategoryRepository.FindByName(ctx, req.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	if existingCategory != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "category already exists")
	}

	// Create new category
	category := &entity.Category{
		Name: req.Name,
	}

	if err := uc.CategoryRepository.Create(tx, category); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to create category: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create category")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create category")
	}

	// Return response
	return converter.CategoryToResponse(category), nil
}

// ListCategories returns a paginated list of categories
func (uc *CategoryUseCase) ListCategories(ctx context.Context, page, size int) ([]*model.CategoryResponse, int64, int64, error) {
	// Ensure default pagination values
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	// Fetch categories with pagination
	var categories []entity.Category
	total, err := uc.CategoryRepository.Paginate(ctx, uc.DB, page, size, &categories)
	if err != nil {
		uc.Log.Error("failed to list categories: ", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, "failed to list categories")
	}

	// Convert entities -> response DTOs via converter
	response := converter.CategoriesToResponse(categories)

	// Calculate total pages
	totalPages := (total + int64(size) - 1) / int64(size)

	// Return result
	return response, total, totalPages, nil
}

func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, id int, req *model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	//  Start transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//  Find existing category
	var category entity.Category
	if err := tx.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "category not found")
		}
		uc.Log.Error("failed to find category: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to find category")
	}

	//  Check duplicate name (exclude itself)
	existingCategory, err := uc.CategoryRepository.FindByName(ctx, req.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	if existingCategory != nil && existingCategory.ID != category.ID {
		return nil, fiber.NewError(fiber.StatusConflict, "category name already exists")
	}

	//  Update fields
	category.Name = req.Name

	if err := uc.CategoryRepository.Update(tx, &category); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to update category: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update category")
	}

	//  Commit transaction
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update category")
	}

	//  Convert entity → response
	return converter.CategoryToResponseModel(&category), nil
}

func (uc *CategoryUseCase) DeleteCategory(ctx context.Context, id int) error {
	// Start transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find existing category
	var category entity.Category
	if err := uc.CategoryRepository.FindById(tx, &category, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "category not found")
		}
		uc.Log.Error("failed to find category: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to find category")
	}

	// Delete category
	if err := uc.CategoryRepository.Delete(tx, &category); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to delete category: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete category")
	}

	// Commit
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete category")
	}

	return nil
}
