package app

import (
	"net/http"

	"github.com/Adedunmol/mycart/internal/routes"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

func Run(logger *httplog.Logger) {

	// Initializers()

	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		util.RespondWithJSON(w, http.StatusOK, "Hello, world")
	})

	routes.SetupRoutes(r)

	http.ListenAndServe(":5000", r)
}
