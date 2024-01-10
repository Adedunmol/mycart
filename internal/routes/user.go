package routes

import "github.com/go-chi/chi/v5"

func UsersRoutes(r *chi.Mux) {

	userRouter := chi.NewRouter()

	r.Mount("/users", userRouter)
}
