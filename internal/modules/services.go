package modules

import (
	"pet-store/internal/infrastructure/component"
	aservice "pet-store/internal/modules/auth/service"
	petservice "pet-store/internal/modules/pet/service"
	uservice "pet-store/internal/modules/user/service"
	oservice "pet-store/internal/modules/order/service"
	"pet-store/internal/storages"
)

type Services struct {
	User uservice.Userer
	Auth aservice.Auther
	Pet  petservice.Peter
	Order oservice.Orderer
}

func NewServices(storages *storages.Storages, components *component.Components) *Services {
	userService := uservice.NewUserService(storages.User, components.Logger)
	return &Services{
		User: userService,
		Auth: aservice.NewAuth(userService, components),
		Pet:  petservice.NewPetService(storages.Pet, components.Logger),
		Order: oservice.NewOrderService(storages.Order, components.Logger),
	}
}
