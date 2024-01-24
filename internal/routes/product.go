package routes

import (
	"github.com/Adedunmol/mycart/internal/services"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
)

func ProductsRoutes(r *chi.Mux) {

	productRouter := chi.NewRouter()

	r.Use(util.AuthMiddleware)

	productRouter.Post("/", services.CreateUserHandler)

	r.Mount("/products", productRouter)
}
