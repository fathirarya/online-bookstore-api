package model

import "time"

// User Response
type UserResponse struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
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

type AuthResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresAt time.Time    `json:"expires_at"`
}

// Logout User
type LogoutUserRequest struct {
	ID int `json:"id" validate:"required"`
}

// Get User
type GetUserRequest struct {
	ID int `json:"id" validate:"required"`
}
