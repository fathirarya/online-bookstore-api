package converter

import (
	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/model"
)

func BookToResponse(book *entity.Book) *model.BookResponse {
	return &model.BookResponse{
		ID:           book.ID,
		Title:        book.Title,
		Author:       book.Author,
		Price:        book.Price,
		Year:         book.Year,
		CategoryID:   book.CategoryID,
		CategoryName: book.Category.Name,
		ImageURL:     book.ImageBase64,
	}
}

func BooksToResponse(books []entity.Book) []*model.BookResponse {
	responses := make([]*model.BookResponse, len(books))
	for i, b := range books {
		responses[i] = &model.BookResponse{
			ID:           b.ID,
			Title:        b.Title,
			Author:       b.Author,
			Price:        b.Price,
			Year:         b.Year,
			CategoryID:   b.CategoryID,
			CategoryName: b.Category.Name,
			ImageURL:     b.ImageBase64,
		}
	}
	return responses
}
