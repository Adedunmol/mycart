package routes

import (
	"github.com/Adedunmol/mycart/internal/services"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
)

func OrdersRoutes(r *chi.Mux) {

	orderRouter := chi.NewRouter()

	orderRouter.Use(util.AuthMiddleware)

	orderRouter.Get("/", services.CreateOrderHandler)

	r.Mount("/orders", orderRouter)
}
