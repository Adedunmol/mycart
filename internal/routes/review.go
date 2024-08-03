package routes

import (
	"github.com/Adedunmol/mycart/internal/services"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
)

func ReviewsRoutes(r *chi.Mux) {

	reviewRouter := chi.NewRouter()

	reviewRouter.Use(util.AuthMiddleware)

	reviewRouter.Post("/{id}", services.CreateReviewHandler)
	reviewRouter.Get("/{review_id}", services.GetReviewHandler)

	r.Mount("/reviews", reviewRouter)
}
