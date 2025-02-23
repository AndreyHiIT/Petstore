package storage

import (
	"context"
	"pet-store/internal/models"
)

type Peter interface {
	AddPet(ctx context.Context, pet *models.Pet) error
	UpdatePet(ctx context.Context, pet models.Pet) error
	FindPetbyStatus(ctx context.Context, statuses []string) ([]models.Pet, error)
	FindPetbyID(ctx context.Context, petID int) (models.Pet, error)
	UpdatePetForm(ctx context.Context, name, status string, petID int) error
}
