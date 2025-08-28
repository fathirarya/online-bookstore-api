package converter

import (
	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/model"
)

// CategoryToResponse converts entity.Category → model.CreateCategoryResponse (used after create)
func CategoryToResponse(category *entity.Category) *model.CreateCategoryResponse {
	return &model.CreateCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}

// CategoriesToResponse converts []entity.Category → []*model.CategoryResponse (used for list)
func CategoriesToResponse(categories []entity.Category) []*model.CategoryResponse {
	res := make([]*model.CategoryResponse, len(categories))
	for i, c := range categories {
		res[i] = &model.CategoryResponse{
			ID:   c.ID,
			Name: c.Name,
		}
	}
	return res
}

func CategoryToResponseModel(c *entity.Category) *model.CategoryResponse {
	return &model.CategoryResponse{
		ID:   c.ID,
		Name: c.Name,
	}
}
