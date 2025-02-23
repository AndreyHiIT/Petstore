package service

import (
	"context"
	"pet-store/internal/models"
)

type Userer interface {
	Create(ctx context.Context, in UserCreateIn) UserCreateOut
	GetByEmail(ctx context.Context, in GetByEmailIn) UserOut
	GetByUsername(ctx context.Context, username string) UserOut
	UpdateUser(ctx context.Context, userdata UpdateUserRequest) UpdateUserResponse
}

type UpdateUserRequest struct {
	UserName  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

type UserCreateIn struct {
	UserName  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserCreateOut struct {
	UserID    int `json:"user_id"`
	ErrorCode int `json:"error_code"`
}

type GetByEmailIn struct {
	Email string `json:"email"`
}

type UserOut struct {
	User      *models.User `json:"user"`
	ErrorCode int          `json:"error_code"`
}

type UpdateUserResponse struct {
	Success   bool `json:"success"`
	ErrorCode int  `json:"error_code,omitempty"`
}
