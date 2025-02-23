package controller

import (
	"net/http"
	"pet-store/internal/infrastructure/component"
	"pet-store/internal/infrastructure/errors"
	"pet-store/internal/infrastructure/responder"
	"pet-store/internal/modules/pet/service"

	"github.com/go-chi/chi"
	"github.com/ptflp/godecoder"
	_ "github.com/vektra/mockery"
)

//go:generate mockery --name=Peter

type Peter interface {
	AddPet(w http.ResponseWriter, r *http.Request)
	UpdatePet(w http.ResponseWriter, r *http.Request)
	FindPetbyStatus(w http.ResponseWriter, r *http.Request)
	FindPetbyID(w http.ResponseWriter, r *http.Request)
	UpdatePetForm(w http.ResponseWriter, r *http.Request)
}

type Pet struct {
	service service.Peter
	responder.Responder
	godecoder.Decoder
}

func NewPet(service service.Peter, components *component.Components) Peter {
	return &Pet{service: service, Responder: components.Responder, Decoder: components.Decoder}
}

// @Summary Update Pet
// @Security ApiKeyAuth
// @Tags pet
// @Description Updates a pet in the store with form data
// @ID UpdatePet
// @Accept  multipart/form-data
// @Produce  json
// @Param petId path string true "Pet ID to update"
// @Param name formData string false "Pet name"
// @Param status formData string false "Pet status"
// @Success 200 {object} SuccessRequest "Successfully updated pet"
// @Failure 400 {object} PetAddResponseErr "Error updating pet"
// @Router /pet/{petId} [post]
func (p *Pet) UpdatePetForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		p.ErrorBadRequest(w, err)
		return
	}

	name := r.FormValue("name")
	status := r.FormValue("status")
	reqID := chi.URLParam(r, "petId")

	out := p.service.UpdatePetForm(r.Context(), name, status, reqID)
	if out.ErrorCode != errors.NoError {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "error update pet"
		p.OutputJSON(w, PetAddResponseErr{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data: Data{
				Message: msg,
			},
		})
		return
	}

	p.OutputJSON(w, SuccessRequest{
		Success: true,
	})
}

// @Summary Find Pet by ID
// @Security ApiKeyAuth
// @Tags pet
// @Description Find a pet by its ID
// @ID FindPetbyID
// @Produce  json
// @Param petId path string true "Pet ID to find"
// @Success 200 {object} FindPetbyIDResponse "Successfully found pet"
// @Router /pet/{petId} [get]
func (p *Pet) FindPetbyID(w http.ResponseWriter, r *http.Request) {
	req := chi.URLParam(r, "petId")

	out := p.service.FindPetbyID(r.Context(), req)
	if out.ErrorCode != errors.NoError {
		var msg string
		if out.ErrorCode == errors.PetServiceErrPetNotFound {
			w.WriteHeader(http.StatusNotFound)
			msg = "no pet found with the provided id"
		} else {
			msg = "Find pet by ID error"
			w.WriteHeader(http.StatusInternalServerError)
		}
		p.OutputJSON(w, PetAddResponseErr{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data: Data{
				Message: msg,
			},
		})
		return
	}

	p.OutputJSON(w, FindPetbyIDResponse{
		Success: true,
		Result:  out.Pet,
	})
}

// FindPetbyStatus godoc
// @Summary Find Pets by Status
// @Security ApiKeyAuth
// @Tags pet
// @Description Find pets based on their status. The status parameter is required. Available values: available, pending, sold.
// @ID FindPetbyStatus
// @Produce  json
// @Param status query string true "Pet status to find" Enums(available, pending, sold)
// @Success 200 {object} FindPetbyStatusResponse "Successfully found pets"
// @Router /pet/findByStatus [get]
func (p *Pet) FindPetbyStatus(w http.ResponseWriter, r *http.Request) {
	req := r.URL.Query()["status"]
	if len(req) == 0 {
		p.OutputJSON(w, "Status parameter is required")
		return
	}

	out := p.service.FindPetbyStatus(r.Context(), req)
	if out.ErrorCode != errors.NoError {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Find pet by status error"
		p.OutputJSON(w, PetAddResponseErr{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data: Data{
				Message: msg,
			},
		})
		return
	}

	p.OutputJSON(w, FindPetbyStatusResponse{
		Success: true,
		Results: out.Pets,
	})
}

// @Summary Add a new Pet
// @Security ApiKeyAuth
// @Tags pet
// @Description Add a new pet to the database with the provided details.
// @ID AddPet
// @Accept  json
// @Produce  json
// @Param pet body service.PetAddRequest true "Pet to add"
// @Success 200 {object} PetAddResponse "Successfully added pet"
// @Router /pet [post]
func (p *Pet) AddPet(w http.ResponseWriter, r *http.Request) {
	var req service.PetAddRequest
	err := p.Decode(r.Body, &req)
	if err != nil {
		p.ErrorBadRequest(w, err)
		return
	}

	out := p.service.AddPet(r.Context(), req)
	if out.ErrorCode != errors.NoError {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "AddPet error"
		p.OutputJSON(w, PetAddResponseErr{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data: Data{
				Message: msg,
			},
		})
		return
	}

	p.OutputJSON(w, PetAddResponse{
		Success: true,
		Data: Data{
			Message: "pet added to the database",
		},
	})
}

// @Summary Update an existing Pet
// @Security ApiKeyAuth
// @Tags pet
// @Description Update an existing pet's details in the database based on the provided pet ID.
// @ID UpdatePet
// @Accept  json
// @Produce  json
// @Param pet body service.PetUpdateRequest true "Pet data to update"
// @Success 200 {object} PetAddResponse "Successfully updated pet"
// @Router /pet [put]
func (p *Pet) UpdatePet(w http.ResponseWriter, r *http.Request) {
	var req service.PetUpdateRequest
	err := p.Decode(r.Body, &req)
	if err != nil {
		p.ErrorBadRequest(w, err)
		return
	}

	out := p.service.UpdatePet(r.Context(), req)
	if out.ErrorCode != errors.NoError {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Update Pet error"
		p.OutputJSON(w, PetAddResponseErr{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data: Data{
				Message: msg,
			},
		})
		return
	}

	p.OutputJSON(w, PetAddResponse{
		Success: true,
		Data: Data{
			Message: "pet success updated",
		},
	})
}
