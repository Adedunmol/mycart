package routes

import (
	"github.com/Adedunmol/mycart/internal/services"
	"github.com/go-chi/chi/v5"
)

func UsersRoutes(r *chi.Mux) {

	userRouter := chi.NewRouter()

	userRouter.Post("/register", services.CreateUserHandler)
	userRouter.Post("/login", services.LoginUserHandler)

	r.Mount("/users", userRouter)
}
