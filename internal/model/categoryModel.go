package model

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type CreateCategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}
