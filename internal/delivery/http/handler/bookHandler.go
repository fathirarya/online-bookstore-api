package handler

import (
	"mime/multipart"
	"strconv"

	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/usecase"
	"github.com/fathirarya/online-bookstore-api/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type BookHandler struct {
	Log      *logrus.Logger
	UseCase  *usecase.BookUseCase
	Validate *validator.Validate
}

func NewBookHandler(useCase *usecase.BookUseCase, logger *logrus.Logger, validate *validator.Validate) *BookHandler {
	return &BookHandler{
		Log:      logger,
		UseCase:  useCase,
		Validate: validate,
	}
}

func (h *BookHandler) Create(ctx *fiber.Ctx) error {
	// Parse non-file fields
	var req model.CreateBookRequest
	req.Title = ctx.FormValue("title")
	req.Author = ctx.FormValue("author")
	req.Price, _ = strconv.ParseFloat(ctx.FormValue("price"), 64)
	req.Year, _ = strconv.Atoi(ctx.FormValue("year"))
	req.CategoryID, _ = strconv.Atoi(ctx.FormValue("category_id"))

	// Validate input fields
	if err := h.Validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "validation failed",
			Errors:  utils.TranslateValidationErrors(err),
		})
	}

	// Handle file upload
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "image file is required",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "failed to open image file",
		})
	}
	defer file.Close()

	// Convert file to base64
	imageBase64, err := utils.FileToBase64(file)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "failed to encode image",
		})
	}

	// Call usecase
	response, err := h.UseCase.CreateBook(ctx.Context(), &req, imageBase64)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	// Success response
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.BookResponse]{
		Data: response,
	})
}

func (h *BookHandler) List(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	size := ctx.QueryInt("size", 10)

	// Call usecase
	books, pageNum, pageSize, totalItems, totalPages, err := h.UseCase.ListBooks(ctx.Context(), page, size)
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

	// Success response with paging
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]*model.BookResponse]{
		Page:       pageNum,
		Size:       pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Data:       books,
	})
}

func (h *BookHandler) GetByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid book id",
		})
	}

	book, err := h.UseCase.GetBookByID(ctx.Context(), id)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.BookResponse]{
		Data:    book,
		Message: "success get book detail",
	})
}

func (h *BookHandler) Update(ctx *fiber.Ctx) error {
	//  Parse ID dari param
	id, err := ctx.ParamsInt("id")
	if err != nil || id <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Message: "invalid book id",
		})
	}

	// Parse form-data ke struct
	var req model.UpdateBookRequest
	req.Title = ctx.FormValue("title")
	req.Author = ctx.FormValue("author")

	// parse optional price
	if priceStr := ctx.FormValue("price"); priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil || price <= 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
				Message: "invalid price value",
			})
		}
		req.Price = price
	}

	// parse optional year
	if yearStr := ctx.FormValue("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
				Message: "invalid year value",
			})
		}
		req.Year = year
	}

	// parse optional category
	if categoryStr := ctx.FormValue("category_id"); categoryStr != "" {
		categoryID, err := strconv.Atoi(categoryStr)
		if err != nil || categoryID <= 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
				Message: "invalid category_id value",
			})
		}
		req.CategoryID = categoryID
	}

	// Ambil file image (optional)
	var file multipart.File
	if fileHeader, err := ctx.FormFile("image"); err == nil {
		f, err := fileHeader.Open()
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
				Message: "failed to open image file",
			})
		}
		defer f.Close()
		file = f
	}

	//  Panggil usecase
	response, err := h.UseCase.UpdateBook(ctx.Context(), id, &req, file)
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

	//  Success response
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.BookResponse]{
		Data: response,
	})
}

func (h *BookHandler) Delete(ctx *fiber.Ctx) error {
	//  Parse ID dari path param
	id, err := ctx.ParamsInt("id")
	if err != nil || id <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Message: "invalid book id",
		})
	}

	//  Panggil usecase untuk delete
	if err := h.UseCase.DeleteBook(ctx.Context(), id); err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return ctx.Status(fiberErr.Code).JSON(model.WebResponse[any]{
				Message: fiberErr.Message,
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Message: "internal server error",
		})
	}

	//  Response sukses
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[any]{
		Message: "book deleted successfully",
	})
}

func (h *BookHandler) GetTotalBooks(ctx *fiber.Ctx) error {
	// 1️⃣ Panggil UseCase
	response, err := h.UseCase.GetTotalBooks(ctx.Context())
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

	// 2️⃣ Response sukses
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.BookStatsResponse]{
		Data:    response,
		Message: "success get total books",
	})
}

func (h *BookHandler) GetBookPriceStats(ctx *fiber.Ctx) error {
	//  Panggil UseCase
	stats, err := h.UseCase.GetBookPriceStats(ctx.Context())
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

	//  Response sukses
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.BookPriceStatsResponse]{
		Data:    stats,
		Message: "success get book price stats",
	})
}
