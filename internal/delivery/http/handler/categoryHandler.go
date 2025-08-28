package handler

import (
	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/usecase"
	"github.com/fathirarya/online-bookstore-api/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CategoryHandler struct {
	Log      *logrus.Logger
	UseCase  *usecase.CategoryUseCase
	Validate *validator.Validate
}

func NewCategoryHandler(useCase *usecase.CategoryUseCase, logger *logrus.Logger, validate *validator.Validate) *CategoryHandler {
	return &CategoryHandler{
		Log:      logger,
		UseCase:  useCase,
		Validate: validate,
	}
}

func (h *CategoryHandler) Create(ctx *fiber.Ctx) error {
	var request model.CreateCategoryRequest

	// 1️⃣ Parse body
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid request body",
		})
	}

	// 2️⃣ Validate struct fields
	if err := h.Validate.Struct(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "validation failed",
			Errors:  utils.TranslateValidationErrors(err), // helper to convert field errors
		})
	}

	// 3️⃣ Call UseCase
	response, err := h.UseCase.CreateCategory(ctx.Context(), &request)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			// Case: category already exists
			if fiberErr.Code == fiber.StatusConflict && fiberErr.Message == "category already exists" {
				return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
					Message: "validation failed",
					Errors:  map[string]string{"name": "category already exists"},
				})
			}

			// Other known errors
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}

		// Unknown error
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	// 4️⃣ Success response
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.CreateCategoryResponse]{
		Data: response,
	})
}

func (h *CategoryHandler) List(ctx *fiber.Ctx) error {
	// Parse query params with default values
	page := ctx.QueryInt("page", 1)
	size := ctx.QueryInt("size", 10)

	// Call UseCase
	data, total, totalPages, err := h.UseCase.ListCategories(ctx.Context(), page, size)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return ctx.Status(fiberErr.Code).JSON(model.WebResponse[any]{
				Message: fiberErr.Message,
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Message: "internal server error",
		})
	}

	// Return success response with pagination metadata
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]*model.CategoryResponse]{
		Page:       page,
		Size:       size,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       data,
	})
}

func (h *CategoryHandler) Update(ctx *fiber.Ctx) error {
	// Parse and validate path param
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid category id",
		})
	}

	// Parse request body
	var request model.UpdateCategoryRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid request body",
		})
	}

	// Validate request struct fields
	if err := h.Validate.Struct(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "validation failed",
			Errors:  utils.TranslateValidationErrors(err),
		})
	}

	// Call UseCase
	response, err := h.UseCase.UpdateCategory(ctx.Context(), id, &request)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			// Case: category not found
			if fiberErr.Code == fiber.StatusNotFound {
				return ctx.Status(fiber.StatusNotFound).JSON(model.ValidationErrorResponse{
					Message: "category not found",
				})
			}
			// Case: duplicate category name
			if fiberErr.Code == fiber.StatusConflict {
				return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
					Message: "validation failed",
					Errors:  map[string]string{"name": "category name already exists"},
				})
			}
			// Other known errors
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}

		// Unknown error
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	// Success response
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.CategoryResponse]{
		Data: response,
	})
}

func (h *CategoryHandler) Delete(ctx *fiber.Ctx) error {
	// Parse and validate path param
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Message: "invalid category id",
		})
	}

	// Call UseCase
	if err := h.UseCase.DeleteCategory(ctx.Context(), id); err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			if fiberErr.Code == fiber.StatusNotFound {
				return ctx.Status(fiber.StatusNotFound).JSON(model.WebResponse[any]{
					Message: "category not found",
				})
			}
			return ctx.Status(fiberErr.Code).JSON(model.WebResponse[any]{
				Message: fiberErr.Message,
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Message: "internal server error",
		})
	}

	// Success response
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[any]{
		Message: "category deleted successfully",
	})
}
