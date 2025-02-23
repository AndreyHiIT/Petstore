package service

import (
	"context"
	"pet-store/internal/infrastructure/errors"
	"pet-store/internal/models"
	"pet-store/internal/modules/pet/storage"
	"strconv"

	"go.uber.org/zap"
)

type PetService struct {
	storage storage.Peter
	logger  *zap.Logger
}

func NewPetService(storage storage.Peter, logger *zap.Logger) *PetService {
	return &PetService{storage: storage, logger: logger}
}

func (p *PetService) UpdatePetForm(ctx context.Context, name, status, reqID string) RequestOut {
	if status == "" && reqID == "" {
		return RequestOut{
			Status:    false,
			ErrorCode: errors.PetServiceUpdatePetFormBadReuest,
		}
	}

	petID, err := strconv.Atoi(reqID)
	if err != nil {
		p.logger.Error("Error during conversion:", zap.Error(err))
		return RequestOut{
			Status:    false,
			ErrorCode: errors.UpdatePetFormErrorDuringConversion,
		}
	}

	err = p.storage.UpdatePetForm(ctx, name, status, petID)
	if err != nil {
		p.logger.Error("Error:", zap.Error(err))
		return RequestOut{
			Status:    false,
			ErrorCode: errors.UpdatePetFormError,
		}
	}

	return RequestOut{
		Status: true,
	}
}

func (p *PetService) FindPetbyID(ctx context.Context, strID string) RequestOutWithPet {
	petID, err := strconv.Atoi(strID)
	if err != nil {
		p.logger.Error("Error during conversion:", zap.Error(err))
		return RequestOutWithPet{
			Status:    false,
			ErrorCode: errors.FindPetbyIDErrorDuringConversion,
		}
	}

	pet, err := p.storage.FindPetbyID(ctx, petID)
	if err != nil {
		p.logger.Error("Error FindPetbyID:", zap.Error(err))
		errorcode := errors.PetServiceFindPetbyID
		if err == errors.ErrPetNotFound {
			errorcode = errors.PetServiceErrPetNotFound
		}
		return RequestOutWithPet{
			Status:    false,
			ErrorCode: errorcode,
		}
	}

	return RequestOutWithPet{
		Pet:    pet,
		Status: true,
	}
}

func (p *PetService) FindPetbyStatus(ctx context.Context, statuses []string) RequestOutWithPets {
	pets, err := p.storage.FindPetbyStatus(ctx, statuses)
	if err != nil {
		p.logger.Error("Error FindPetbyStatus:", zap.Error(err))
		return RequestOutWithPets{
			Status:    false,
			ErrorCode: errors.PetServiceFindPetbyStatus,
		}
	}

	return RequestOutWithPets{
		Pets:   pets,
		Status: true,
	}
}

func (p *PetService) AddPet(ctx context.Context, pet PetAddRequest) RequestOut {
	photourls := make([]string, 0, len(pet.PhotoUrls))
	if pet.PhotoUrls != nil {
		photourls = append(photourls, pet.PhotoUrls...)
	}

	tags := make([]models.Tag, 0, len(pet.Tags))
	if pet.Tags != nil {
		for i := range pet.Tags {
			tags = append(tags, models.Tag{Name: pet.Tags[i].Name})
		}
	}

	petStorage := models.Pet{
		Category:  models.Category{Name: pet.Category.Name},
		Name:      pet.Name,
		PhotoUrls: photourls,
		Tags:      tags,
		Status:    pet.Status,
	}
	err := p.storage.AddPet(ctx, &petStorage)
	if err != nil {
		p.logger.Error("Error AddPet:", zap.Error(err))
		return RequestOut{
			Status:    false,
			ErrorCode: errors.AddPetErr,
		}
	}

	return RequestOut{
		Status: true,
	}
}

func (p *PetService) UpdatePet(ctx context.Context, pet PetUpdateRequest) RequestOut {
	photourls := make([]string, 0, len(pet.PhotoUrls))
	if pet.PhotoUrls != nil {
		photourls = append(photourls, pet.PhotoUrls...)
	}

	tags := make([]models.Tag, 0, len(pet.Tags))
	if pet.Tags != nil {
		for i := range pet.Tags {
			tags = append(tags, models.Tag{Name: pet.Tags[i].Name})
		}
	}

	petStorage := models.Pet{
		ID:        pet.ID,
		Category:  models.Category{Name: pet.Category.Name},
		Name:      pet.Name,
		PhotoUrls: photourls,
		Tags:      tags,
		Status:    pet.Status,
	}

	err := p.storage.UpdatePet(ctx, petStorage)
	if err != nil {
		p.logger.Error("Error UpdatePet:", zap.Error(err))
		return RequestOut{
			Status:    false,
			ErrorCode: errors.PetServiceUpdateErr,
		}
	}

	return RequestOut{
		Status: true,
	}
}
