package router

import (
	"pet-store/internal/infrastructure/component"
	"pet-store/internal/modules"
	"pet-store/internal/router"

	"github.com/go-chi/chi"
)

func NewRouter(controllers *modules.Controllers, components *component.Components) *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/", router.NewApiRouter(controllers, components))
	return r
}
