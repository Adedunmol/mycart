package app

import (
	"log"
	"net/http"
	"sync"

	"github.com/Adedunmol/mycart/internal/config"
	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/Adedunmol/mycart/internal/routes"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var dbConnOnce sync.Once

func init() {
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	dbConnOnce.Do(func() {
		err = database.InitDB()

		if err != nil {
			log.Panic(err)
		}

		database.InsertRoles()
	})
}

func Run() {
	logger.Logger.Info("setting up routes")

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		util.RespondWithJSON(w, http.StatusOK, "Hello, world")
	})

	routes.SetupRoutes(r)

	http.ListenAndServe(":5000", r)
}
