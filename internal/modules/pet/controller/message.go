package controller

import "pet-store/internal/models"

type PetAddResponseErr struct {
	Success   bool
	ErrorCode int
	Data      Data
}

type PetAddResponse struct {
	Success bool
	Data    Data
}

type Data struct {
	Message string
}

type FindPetbyStatusResponse struct {
	Success bool
	Results []models.Pet
}

type FindPetbyIDResponse struct {
	Success bool
	Result  models.Pet
}

type SuccessRequest struct {
	Success bool
}
