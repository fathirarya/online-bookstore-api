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
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	BookRepository     *repository.BookRepository
	CategoryRepository *repository.CategoryRepository
}

func NewBookUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	bookRepository *repository.BookRepository, categoryRepository *repository.CategoryRepository) *BookUseCase {
	return &BookUseCase{
		DB:                 db,
		Log:                logger,
		Validate:           validate,
		BookRepository:     bookRepository,
		CategoryRepository: categoryRepository,
	}
}

func (uc *BookUseCase) CreateBook(ctx context.Context, req *model.CreateBookRequest, file multipart.File) (*model.BookResponse, error) {
	// Start transaction
	tx := uc.DB.Begin()
	if tx.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to start transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validate input (kecuali image karena dari file)
	if err := uc.Validate.StructExcept(req, "ImageBase64"); err != nil {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusBadRequest, "validation failed, please check your input")
	}

	// Convert file ke base64
	imageBase64, err := utils.FileToBase64(file)
	if err != nil {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to process image")
	}

	// Cek apakah buku sudah ada berdasarkan title
	existingBook, err := uc.BookRepository.FindByTitle(ctx, req.Title)
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	if existingBook != nil {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusConflict, "book already exists")
	}

	// Buat entitas Book baru
	book := &entity.Book{
		Title:       req.Title,
		Author:      req.Author,
		Price:       req.Price,
		Year:        req.Year,
		CategoryID:  req.CategoryID,
		ImageBase64: imageBase64,
	}

	// Simpan ke database dengan transaction
	if err := uc.BookRepository.Create(tx, book); err != nil {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create book")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create book")
	}

	// Return response (tanpa log berlebihan)
	return converter.BookToResponse(book), nil
}

func (uc *BookUseCase) ListBooks(ctx context.Context, page, size int) ([]*model.BookResponse, *model.PageMetadata, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	var books []entity.Book
	total, err := uc.BookRepository.Paginate(ctx, uc.DB, page, size, &books)
	if err != nil {
		uc.Log.Error("failed to list categories: ", err)
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, "failed to list categories")
	}

	response := make([]*model.BookResponse, len(books))
	for i, b := range books {
		response[i] = &model.BookResponse{
			ID:         b.ID,
			Title:      b.Title,
			Author:     b.Author,
			Price:      b.Price,
			Year:       b.Year,
			CategoryID: b.CategoryID,
			ImageURL:   b.ImageBase64,
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

func (uc *BookUseCase) GetBookByID(ctx context.Context, id int) (*model.BookResponse, error) {
	book := &entity.Book{}

	tx := uc.DB.Begin()
	if tx.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to start transaction")
	}

	// Query ke repository
	err := uc.BookRepository.FindById(tx, book, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "book not found")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to get book")
	}

	return converter.BookToResponse(book), nil
}

func (uc *BookUseCase) UpdateBook(ctx context.Context, id int, req *model.UpdateBookRequest, file multipart.File) (*model.BookResponse, error) {
	// Start transaction
	tx := uc.DB.Begin()
	if tx.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to start transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validate input (kecuali image karena dari file)
	if err := uc.Validate.StructExcept(req, "ImageBase64"); err != nil {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusBadRequest, "validation failed, please check your input")
	}

	// Ambil data lama
	var book entity.Book
	if err := tx.First(&book, id).Error; err != nil {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusNotFound, "book not found")
	}

	// Validasi per field satu-satu

	// Cek Title
	if req.Title != "" && req.Title != book.Title {
		existingBook, err := uc.BookRepository.FindByTitle(ctx, req.Title)
		if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}
		if existingBook != nil && existingBook.ID != id {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusConflict, "book title already exists")
		}
		book.Title = req.Title
	}

	// Cek Author
	if req.Author != "" && req.Author != book.Author {
		existingBook, err := uc.BookRepository.FindByAuthor(ctx, req.Author)
		if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}
		if existingBook != nil && existingBook.ID != id {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusConflict, "book author already exists")
		}
		book.Author = req.Author
	}

	// Cek Price
	if req.Price != 0 && req.Price != book.Price {
		existingBook, err := uc.BookRepository.FindByPrice(ctx, req.Price)
		if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}
		if existingBook != nil && existingBook.ID != id {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusConflict, "book price already exists")
		}
		book.Price = req.Price
	}

	// Cek Year
	if req.Year != 0 && req.Year != book.Year {
		existingBook, err := uc.BookRepository.FindByYear(ctx, req.Year)
		if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}
		if existingBook != nil && existingBook.ID != id {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusConflict, "book year already exists")
		}
		book.Year = req.Year
	}

	// Cek CategoryID
	if req.CategoryID != 0 && req.CategoryID != book.CategoryID {
		existingBook, err := uc.BookRepository.FindByCategoryID(ctx, req.CategoryID)
		if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}
		if existingBook != nil && existingBook.ID != id {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusConflict, "category already assigned to another book")
		}
		book.CategoryID = req.CategoryID
	}

	// Update Image jika ada file baru
	if file != nil {
		imageBase64, err := utils.FileToBase64(file)
		if err != nil {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to process image")
		}
		existingBook, err := uc.BookRepository.FindByImage(ctx, imageBase64)
		if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}
		if existingBook != nil && existingBook.ID != id {
			tx.Rollback()
			return nil, fiber.NewError(fiber.StatusConflict, "image already used by another book")
		}
		book.ImageBase64 = imageBase64
	}

	// Simpan perubahan
	if err := uc.BookRepository.Update(tx, &book); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update book")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update book")
	}

	// Return response
	return converter.BookToResponse(&book), nil
}

func (uc *BookUseCase) DeleteBook(ctx context.Context, id int) error {
	// Mulai transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Ambil data lama
	var book entity.Book
	if err := tx.First(&book, id).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusNotFound, "book not found")
	}

	// Hapus data
	if err := uc.BookRepository.Delete(tx, &book); err != nil {
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
	// Mulai transaksi (read-only)
	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Ambil statistik harga dari repository
	stats, err := uc.BookRepository.GetPriceStats()
	if err != nil {
		tx.Rollback()
		uc.Log.Error("failed to get book price stats: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to get book price stats")
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to get book price stats")
	}

	return stats, nil
}
