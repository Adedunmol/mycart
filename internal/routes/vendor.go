package routes

import (
	"github.com/Adedunmol/mycart/internal/services"
	"github.com/go-chi/chi/v5"
)

func VendorsRoutes(r *chi.Mux) {

	vendorsRouter := chi.NewRouter()

	vendorsRouter.Post("/register", services.CreateUserHandler)

	r.Mount("/vendors", vendorsRouter)
}
