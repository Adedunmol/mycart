package routes

import (
	"github.com/go-chi/chi/v5"
)

func CartsRoutes(r *chi.Mux) {

	cartRouter := chi.NewRouter()

	r.Mount("/carts", cartRouter)
}
