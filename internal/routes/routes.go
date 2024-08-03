package routes

import "github.com/go-chi/chi/v5"

func SetupRoutes(r *chi.Mux) {

	UsersRoutes(r)
	VendorsRoutes(r)
	ProductsRoutes(r)
	CartsRoutes(r)
	OrdersRoutes(r)
	ReviewsRoutes(r)
}
