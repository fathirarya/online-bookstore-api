package converter

import (
	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/model"
)

func BookToResponse(book *entity.Book) *model.BookResponse {
	return &model.BookResponse{
		ID:         book.ID,
		Title:      book.Title,
		Author:     book.Author,
		Price:      book.Price,
		Year:       book.Year,
		CategoryID: book.CategoryID,
		ImageURL:   book.ImageBase64,
	}
}
