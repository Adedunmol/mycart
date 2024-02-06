package routes

import (
	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/services"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
)

func ProductsRoutes(r *chi.Mux) {

	productRouter := chi.NewRouter()

	// r.Use(util.AuthMiddleware)

	productRouter.Use(util.AuthMiddleware)

	productRouter.Get("/", services.GetAllProductsHandler)
	productRouter.Get("/{id}", services.GetProductHandler)

	productRouter.Group(func(r chi.Router) {
		r.Use(util.RoleAuthorization(database.CREATE_PRODUCT))

		r.Post("/", services.CreateProductHandler)
		r.Delete("/{id}", services.DeleteProductHandler)
	})

	productRouter.Group(func(r chi.Router) {
		r.Use(util.RoleAuthorization(database.MODIFY_PRODUCT))

		r.Patch("/{id}", services.UpdateProductHandler)
	})

	r.Mount("/products", productRouter)
}
