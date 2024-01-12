package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/routes"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

func Initializers() (database.DbInstance, util.Config) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	db, err := database.InitDB(config)

	if err != nil {
		log.Panic(err)
	}

	return db, config
}

func Run(logger *httplog.Logger) {

	db, _ := Initializers()

	fmt.Println("db: ", db)

	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		util.RespondWithJSON(w, http.StatusOK, "Hello, world")
	})

	routes.SetupRoutes(r)

	http.ListenAndServe(":5000", r)
}
