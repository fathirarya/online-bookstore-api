package converter

import (
	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/model"
)

func CategoryToResponse(user *entity.Category) *model.CreateCategoryResponse {
	return &model.CreateCategoryResponse{
		ID:   user.ID,
		Name: user.Name,
	}
}

func CategoriesToResponse(categories []*entity.Category) []*model.CategoryResponse {
	var categoryResponses []*model.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, &model.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}
	return categoryResponses
}
