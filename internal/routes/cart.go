package routes

import (
	"github.com/Adedunmol/mycart/internal/services"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
)

func CartsRoutes(r *chi.Mux) {

	cartRouter := chi.NewRouter()

	cartRouter.Use(util.AuthMiddleware)

	cartRouter.Get("/", services.GetCartHandler)
	cartRouter.Post("/add-item", services.AddToRedisCartHandler)
	cartRouter.Post("/remove-item", services.RemoveFromRedisCartHandler)

	r.Mount("/carts", cartRouter)
}
