package app

import (
	"log"
	"net/http"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/routes"
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

	Initializers()

	r := chi.NewRouter()

	routes.SetupRoutes(r)

	r.Use(httplog.RequestLogger(logger))

	http.ListenAndServe("127.0.0.1:5000", r)
}
