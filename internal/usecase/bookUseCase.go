package usecase

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/model/converter"
	"github.com/fathirarya/online-bookstore-api/internal/repository"
	"github.com/fathirarya/online-bookstore-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	BookRepository     *repository.BookRepository
	CategoryRepository *repository.CategoryRepository
}

func NewBookUseCase(db *gorm.DB, logger *logrus.Logger, bookRepository *repository.BookRepository,
	categoryRepository *repository.CategoryRepository) *BookUseCase {
	return &BookUseCase{
		DB:                 db,
		Log:                logger,
		BookRepository:     bookRepository,
		CategoryRepository: categoryRepository,
	}
}

func (uc *BookUseCase) CreateBook(ctx context.Context, req *model.CreateBookRequest, imageBase64 string) (*model.BookResponse, error) {
	// Start transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check duplicate title
	existingBook, err := uc.BookRepository.FindByTitle(ctx, req.Title)
	if err != nil && err != gorm.ErrRecordNotFound {
		uc.Log.Error("failed to check existing book: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to check book")
	}
	if existingBook != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "book already exists")
	}

	// 3️⃣ Build entity
	book := &entity.Book{
		Title:       req.Title,
		Author:      req.Author,
		Price:       req.Price,
		Year:        req.Year,
		CategoryID:  req.CategoryID,
		ImageBase64: imageBase64,
	}

	// Persist book
	if err := uc.BookRepository.Create(tx, book); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to create book: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create book")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create book")
	}

	// Fetch full book with category preload (supaya CategoryName ikut keisi)
	fullBook, err := uc.BookRepository.FindByID(ctx, uc.DB, book.ID)
	if err != nil {
		uc.Log.Error("failed to fetch full book after create: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch book")
	}

	// Convert to response
	return converter.BookToResponse(fullBook), nil
}

func (uc *BookUseCase) ListBooks(ctx context.Context, page, size int) ([]*model.BookResponse, int, int, int64, int64, error) {
	// Default pagination
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	// Query books
	var books []entity.Book
	total, err := uc.BookRepository.Paginate(ctx, uc.DB.Preload("Category"), page, size, &books)
	if err != nil {
		uc.Log.Error("failed to list books: ", err)
		return nil, 0, 0, 0, 0, fiber.NewError(fiber.StatusInternalServerError, "failed to list books")
	}

	// Convert to response
	response := converter.BooksToResponse(books)

	// Hitung total pages
	totalPages := (total + int64(size) - 1) / int64(size)

	return response, page, size, total, totalPages, nil
}

func (uc *BookUseCase) GetBookByID(ctx context.Context, id int) (*model.BookResponse, error) {
	var book entity.Book
	if err := uc.BookRepository.FindById(uc.DB.WithContext(ctx), &book, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "book not found")
		}
		uc.Log.Error("failed to get book: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to get book")
	}

	return converter.BookToResponse(&book), nil
}

func (uc *BookUseCase) UpdateBook(ctx context.Context, id int, req *model.UpdateBookRequest, file multipart.File) (*model.BookResponse, error) {
	// Start transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Ambil data lama
	var book entity.Book
	if err := tx.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "book not found")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch book")
	}

	// Cek Title unik (hanya Title yang perlu unik)
	if req.Title != "" && req.Title != book.Title {
		existingBook, err := uc.BookRepository.FindByTitle(ctx, req.Title)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			uc.Log.Error("failed to check title uniqueness: ", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to check book title")
		}
		if existingBook != nil && existingBook.ID != id {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusConflict, "book title already exists")
		}
		book.Title = req.Title
	}

	// Update field lain (tidak perlu cek unik)
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Price > 0 {
		book.Price = req.Price
	}
	if req.Year > 0 {
		book.Year = req.Year
	}
	if req.CategoryID > 0 {
		book.CategoryID = req.CategoryID
	}

	// Update Image jika ada file baru
	if file != nil {
		imageBase64, err := utils.FileToBase64(file)
		if err != nil {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to process image")
		}
		book.ImageBase64 = imageBase64
	}

	// Simpan perubahan
	if err := uc.BookRepository.Update(tx, &book); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to update book: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update book")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update book")
	}

	// Ambil lagi dengan preload Category agar response lengkap
	updatedBook, err := uc.BookRepository.FindByID(ctx, uc.DB, book.ID)
	if err != nil {
		uc.Log.Error("failed to fetch updated book: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch updated book")
	}

	// Return response
	return converter.BookToResponse(updatedBook), nil
}

func (uc *BookUseCase) DeleteBook(ctx context.Context, id int) error {
	// Start transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//  Ambil data lama via repository
	book, err := uc.BookRepository.FindByID(ctx, tx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "book not found")
		}
		uc.Log.Error("failed to fetch book before delete: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch book")
	}

	//  Hapus data
	if err := uc.BookRepository.Delete(tx, book); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to delete book: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete book")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete book")
	}

	return nil
}

func (uc *BookUseCase) GetTotalBooks(ctx context.Context) (*model.BookStatsResponse, error) {
	// Mulai transaction (meskipun read-only)
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Hitung total buku menggunakan repository
	total, err := uc.BookRepository.CountAllBooks()
	if err != nil {
		tx.Rollback()
		uc.Log.Error("failed to count books: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to count books")
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to count books")
	}

	// Kembalikan response
	return &model.BookStatsResponse{
		TotalBooks: int(total),
	}, nil
}

func (uc *BookUseCase) GetBookPriceStats(ctx context.Context) (*model.BookPriceStatsResponse, error) {
	// 1️⃣ Ambil statistik harga dari repository
	stats, err := uc.BookRepository.GetPriceStats(ctx, uc.DB)
	if err != nil {
		uc.Log.Error("failed to get book price stats: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to get book price stats")
	}

	// 2️⃣ Return response
	return stats, nil
}
