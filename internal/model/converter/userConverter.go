package converter

import (
	"time"

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

func AuthToResponse(user *entity.User, token string, expiresAt time.Time) *model.AuthResponse {
	return &model.AuthResponse{
		Token:     token,
		User:      *UserToResponse(user),
		ExpiresAt: expiresAt,
	}
}
