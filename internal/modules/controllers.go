package modules

import (
	"pet-store/internal/infrastructure/component"
	acontroller "pet-store/internal/modules/auth/controller"
	ocontroller "pet-store/internal/modules/order/controller"
	pcontroller "pet-store/internal/modules/pet/controller"
	ucontroller "pet-store/internal/modules/user/controller"
)

type Controllers struct {
	Auth  acontroller.Auther
	User  ucontroller.Userer
	Pet   pcontroller.Peter
	Order ocontroller.Orderer
}

func NewControllers(services *Services, components *component.Components) *Controllers {
	authController := acontroller.NewAuth(services.Auth, components)
	userController := ucontroller.NewUser(services.User, components)
	petController := pcontroller.NewPet(services.Pet, components)
	orderController := ocontroller.NewOrder(services.Order, components)
	return &Controllers{
		Auth:  authController,
		User:  userController,
		Pet:   petController,
		Order: orderController,
	}
}
