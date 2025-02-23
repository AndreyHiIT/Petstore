package router

import (
	"net/http"
	_ "pet-store/docs"
	"pet-store/internal/infrastructure/component"
	"pet-store/internal/modules"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewApiRouter(controllers *modules.Controllers, components *component.Components) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// Swagger route
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Auth routes
	r.Route("/user", func(r chi.Router) {
		authController := controllers.Auth
		userController := controllers.User
		r.Post("/", authController.CreateUser)
		r.Get("/login", authController.Login)
		//r.Post("/logout", authController.Logout)
		//r.Post("/createWithList", authController.CreateWithList)
		r.Get("/{username}", userController.GetUser)
		r.Put("/{username}", userController.UpdateUser)
		//r.Delete("/{username}", c.deleteUser)
		//r.Post("/createWithArray", c.createWithArray)
	})
	r.Group(func(r chi.Router) {
		r.Route("/pet", func(r chi.Router) {
			r.Use(jwtauth.Verifier(components.TokenManager.GetAccessSecret()))
			r.Use(jwtauth.Authenticator)
			petController := controllers.Pet
			//	r.Post("/{petId}/uploadImage", c.uploadImage)
			r.Post("/", petController.AddPet)
			r.Put("/", petController.UpdatePet)
			r.Get("/findByStatus", petController.FindPetbyStatus)
			r.Get("/{petId}", petController.FindPetbyID)
			r.Post("/{petId}", petController.UpdatePetForm)
			//r.Delete("/{petId}", petController.DeletePet)
		})
	})
	r.Route("/store", func(r chi.Router) {
		r.Route("/order", func(r chi.Router) {
			orderController := controllers.Order
			r.Post("/", orderController.CreateOrder)
			r.Get("/{orderId}", orderController.FindOrderByID)
			r.Delete("/{orderId}", orderController.DeleteOrderByID)
		})
		//r.With(jwtauth.Verifier(components.TokenManager.GetAccessSecret()), jwtauth.Authenticator).
		//	Get("/inventory", orderController.Inventory)
	})

	return r
}
