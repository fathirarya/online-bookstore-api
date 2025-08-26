package model

import "time"

// User Response
type UserResponse struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Verify User
type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}

// Register User
type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=255"`
}

// Update User
type UpdateUserRequest struct {
	ID       int    `json:"-" validate:"required"` // INT sesuai entity.User
	Name     string `json:"name,omitempty" validate:"max=100"`
	Password string `json:"password,omitempty" validate:"max=255"`
}

// Login User
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=255"`
}

// Logout User
type LogoutUserRequest struct {
	ID int `json:"id" validate:"required"`
}

// Get User
type GetUserRequest struct {
	ID int `json:"id" validate:"required"`
}
