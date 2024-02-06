package routes

import (
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
)

func CartsRoutes(r *chi.Mux) {

	cartRouter := chi.NewRouter()

	cartRouter.Use(util.AuthMiddleware)

	r.Mount("/carts", cartRouter)
}
