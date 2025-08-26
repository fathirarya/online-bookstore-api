package converter

import (
	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}
