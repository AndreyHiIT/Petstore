package storage

import (
	"context"
	"pet-store/internal/db/adapter"
	"pet-store/internal/models"
)

// PetStorage - хранилище животных
type PetStorage struct {
	adapter *adapter.SQLAdapter
}

// NewUserStorage - конструктор хранилища пользователей
func NewPetStorage(sqlAdapter *adapter.SQLAdapter) *PetStorage {
	return &PetStorage{adapter: sqlAdapter}
}

func (a *PetStorage) AddPet(ctx context.Context, pet *models.Pet) error {
	err := a.adapter.AddPet(ctx, pet)
	return err
}

func (a *PetStorage) UpdatePet(ctx context.Context, pet models.Pet) error {
	return a.adapter.UpdatePet(ctx, pet)
}

func (a *PetStorage) FindPetbyStatus(ctx context.Context, statuses []string) ([]models.Pet, error) {
	return a.adapter.FindPetbyStatus(ctx, statuses)
}

func (a *PetStorage) FindPetbyID(ctx context.Context, petID int) (models.Pet, error) {
	return a.adapter.FindPetbyID(ctx, petID)
}

func (a *PetStorage) UpdatePetForm(ctx context.Context, name, status string, petID int) error {
	return a.adapter.UpdatePetForm(ctx, name, status, petID)
}
