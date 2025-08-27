package handler

import (
	"errors"
	"mime/multipart"
	"strconv"

	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookHandler struct {
	Log     *logrus.Logger
	UseCase *usecase.BookUseCase
}

func NewBookHandler(useCase *usecase.BookUseCase, logger *logrus.Logger) *BookHandler {
	return &BookHandler{
		Log:     logger,
		UseCase: useCase,
	}
}

func (h *BookHandler) Create(ctx *fiber.Ctx) error {
	// ambil form-data
	title := ctx.FormValue("title")
	author := ctx.FormValue("author")
	priceStr := ctx.FormValue("price")
	yearStr := ctx.FormValue("year")
	categoryStr := ctx.FormValue("category_id")

	// validasi input kosong
	if title == "" || author == "" || priceStr == "" || categoryStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "title, author, price, and category_id are required",
		})
	}

	// parse angka
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid price value",
		})
	}

	year := 0
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
				Message: "invalid year value",
			})
		}
	}

	categoryID, err := strconv.Atoi(categoryStr)
	if err != nil || categoryID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid category_id value",
		})
	}

	// ambil file (image)
	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "image file is required",
		})
	}

	// open file
	f, err := file.Open()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "failed to open image file",
		})
	}
	defer f.Close()

	// map ke request struct
	req := &model.CreateBookRequest{
		Title:      title,
		Author:     author,
		Price:      price,
		Year:       year,
		CategoryID: categoryID,
	}

	// panggil usecase
	response, err := h.UseCase.CreateBook(ctx.Context(), req, f)
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

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.BookResponse]{
		Data: response,
	})
}

func (h *BookHandler) List(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	size := ctx.QueryInt("size", 10)

	books, paging, err := h.UseCase.ListBooks(ctx.Context(), page, size)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ValidationErrorResponse{
				Message: "books not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]*model.BookResponse]{
		Paging: paging,
		Data:   books,
	})
}

func (h *BookHandler) GetByID(ctx *fiber.Ctx) error {
	// Ambil parameter id dari URL
	id, err := ctx.ParamsInt("id")
	if err != nil || id <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid book id",
		})
	}

	// Panggil usecase
	book, err := h.UseCase.GetBookByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(model.ValidationErrorResponse{
				Message: "book not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	// Return success response
	return ctx.Status(fiber.StatusOK).JSON(book)
}

func (h *BookHandler) Update(ctx *fiber.Ctx) error {
	// ambil param id
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid book id",
		})
	}

	// ambil form-data
	title := ctx.FormValue("title")
	author := ctx.FormValue("author")
	priceStr := ctx.FormValue("price")
	yearStr := ctx.FormValue("year")
	categoryStr := ctx.FormValue("category_id")

	// parsing number fields (optional)
	var price float64
	if priceStr != "" {
		price, err = strconv.ParseFloat(priceStr, 64)
		if err != nil || price <= 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
				Message: "invalid price value",
			})
		}
	}

	var year int
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
				Message: "invalid year value",
			})
		}
	}

	var categoryID int
	if categoryStr != "" {
		categoryID, err = strconv.Atoi(categoryStr)
		if err != nil || categoryID <= 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
				Message: "invalid category_id value",
			})
		}
	}

	// ambil file image (optional)
	var file multipart.File
	fileHeader, err := ctx.FormFile("image")
	if err == nil {
		f, err := fileHeader.Open()
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
				Message: "failed to open image file",
			})
		}
		defer f.Close()
		file = f
	}

	// map ke request struct
	req := &model.UpdateBookRequest{
		Title:      title,
		Author:     author,
		Price:      price,
		Year:       year,
		CategoryID: categoryID,
	}

	// panggil usecase
	response, err := h.UseCase.UpdateBook(ctx.Context(), id, req, file)
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
		Data: response,
	})
}

func (h *BookHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[string]{
			Errors: "invalid book id",
		})
	}

	if err := h.UseCase.DeleteBook(ctx.Context(), id); err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return ctx.Status(fiberErr.Code).JSON(model.WebResponse[string]{
				Errors: fiberErr.Message,
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[string]{
			Errors: "internal server error",
		})
	}

	// Return response konfirmasi delete berhasil
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[string]{
		Data: "book deleted successfully",
	})
}

func (h *BookHandler) GetTotalBooks(ctx *fiber.Ctx) error {
	// Panggil UseCase untuk menghitung total buku
	response, err := h.UseCase.GetTotalBooks(ctx.Context())
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return ctx.Status(fiberErr.Code).JSON(model.WebResponse[string]{
				Errors: fiberErr.Message,
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[string]{
			Errors: "internal server error",
		})
	}

	// Return response sukses
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.BookStatsResponse]{
		Data: response,
	})
}

func (h *BookHandler) GetBookPriceStats(ctx *fiber.Ctx) error {
	// Panggil UseCase untuk mendapatkan statistik harga buku
	stats, err := h.UseCase.GetBookPriceStats(ctx.Context())
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return ctx.Status(fiberErr.Code).JSON(model.WebResponse[string]{
				Errors: fiberErr.Message,
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[string]{
			Errors: "internal server error",
		})
	}

	// Response sukses
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.BookPriceStatsResponse]{
		Data: stats,
	})
}
