package app

import (
	"log"
	"net/http"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/routes"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

func Initializers() database.DbInstance {
	db, err := database.InitDB()

	if err != nil {
		log.Panic(err)
	}

	return db
}

func Run(logger *httplog.Logger) {

	// Initializers()

	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		util.RespondWithJSON(w, http.StatusOK, "")
	})

	routes.SetupRoutes(r)

	http.ListenAndServe("127.0.0.1:5000", r)
}
