package storage

import (
	"context"
	"pet-store/internal/models"
)

type Userer interface {
	Create(ctx context.Context, u models.UserDTO) (int, error)
	GetByEmail(ctx context.Context, email string) (models.UserDTO, error)
	GetByUsername(ctx context.Context, username string) (models.UserDTO, error)
	UpdateUser(ctx context.Context, u models.UserDTO) error
}
