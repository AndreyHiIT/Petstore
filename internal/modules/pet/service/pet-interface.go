package service

import (
	"context"
	"pet-store/internal/models"
)

//go:generate mockery --name=Peter

type Peter interface {
	AddPet(ctx context.Context, pet PetAddRequest) RequestOut
	UpdatePet(ctx context.Context, pet PetUpdateRequest) RequestOut
	FindPetbyStatus(ctx context.Context, statuses []string) RequestOutWithPets
	FindPetbyID(ctx context.Context, strID string) RequestOutWithPet
	UpdatePetForm(ctx context.Context, name, status, reqID string) RequestOut
}

type PetAddRequest struct {
	Category  Category `json:"category"`
	Name      string   `json:"name" `
	PhotoUrls []string `json:"photourls"`
	Tags      []Tag    `json:"tags"`
	Status    string   `json:"status"`
}

type Category struct {
	Name string `json:"name"`
}

type Tag struct {
	Name string `json:"name"`
}

type RequestOut struct {
	Status    bool
	ErrorCode int
}

type RequestOutWithPet struct {
	Pet       models.Pet
	Status    bool
	ErrorCode int
}

type RequestOutWithPets struct {
	Pets      []models.Pet
	Status    bool
	ErrorCode int
}

type PetUpdateRequest struct {
	ID        int      `json:"id"`
	Category  Category `json:"category"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photourls"`
	Tags      []Tag    `json:"tags"`
	Status    string   `json:"status"`
}

type PetFindResponse struct {
	ID        int      `json:"id"`
	Category  Category `json:"category"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photourls"`
	Tags      []Tag    `json:"tags"`
	Status    string   `json:"status"`
}
