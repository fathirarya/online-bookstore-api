package handler

import (
	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CategoryHandler struct {
	Log     *logrus.Logger
	UseCase *usecase.CategoryUseCase
}

func NewCategoryHandler(useCase *usecase.CategoryUseCase, logger *logrus.Logger) *CategoryHandler {
	return &CategoryHandler{
		Log:     logger,
		UseCase: useCase,
	}
}

func (h *CategoryHandler) Create(ctx *fiber.Ctx) error {
	var request model.CreateCategoryRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid request body",
		})
	}

	response, err := h.UseCase.CreateCategory(ctx.Context(), &request)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			if fiberErr.Code == fiber.StatusConflict && fiberErr.Message == "category already exists" {
				return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
					Message: "validation failed",
					Errors:  map[string]string{"name": "category already exists"},
				})
			}
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.CreateCategoryResponse]{
		Data: response,
	})

}

func (h *CategoryHandler) List(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	size := ctx.QueryInt("size", 10)

	categories, paging, err := h.UseCase.ListCategories(ctx.Context(), page, size)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return ctx.Status(fiberErr.Code).JSON(model.WebResponse[string]{Errors: fiberErr.Message})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[string]{Errors: "internal server error"})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]*model.CategoryResponse]{
		Paging: paging,
		Data:   categories,
	})
}

func (h *CategoryHandler) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid category id",
		})
	}

	var request model.UpdateCategoryRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid request body",
		})
	}

	response, err := h.UseCase.UpdateCategory(ctx.Context(), id, &request)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			if fiberErr.Code == fiber.StatusNotFound && fiberErr.Message == "category not found" {
				return ctx.Status(fiber.StatusNotFound).JSON(model.ValidationErrorResponse{
					Message: "category not found",
				})
			}
			if fiberErr.Code == fiber.StatusConflict && fiberErr.Message == "category name already exists" {
				return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
					Message: "validation failed",
					Errors:  map[string]string{"name": "category name already exists"},
				})
			}
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.CategoryResponse]{
		Data: response,
	})
}

func (h *CategoryHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[string]{
			Errors: "invalid category id",
		})
	}

	if err := h.UseCase.DeleteCategory(ctx.Context(), id); err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return ctx.Status(fiberErr.Code).JSON(model.WebResponse[string]{Errors: fiberErr.Message})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[string]{Errors: "internal server error"})
	}

	// Return response konfirmasi delete berhasil
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[string]{
		Data: "category deleted successfully",
	})
}
