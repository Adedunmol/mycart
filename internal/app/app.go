package app

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/Adedunmol/mycart/internal/config"
	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/Adedunmol/mycart/internal/routes"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	dbConnOnce sync.Once
	Router     *chi.Mux
)

func Init() {

	dbConnOnce.Do(func() {
		var err error
		if config.EnvConfig.Environment == "test" {
			err = database.InitDB(config.EnvConfig.TestDatabaseUrl)
		} else {
			err = database.InitDB(config.EnvConfig.DatabaseUrl)
		}

		if err != nil {
			log.Panic(err)
		}

		database.InsertRoles()
	})
}

func Run() {
	// h := asynqmon.New(asynqmon.Options{
	// 	RootPath:     "/monitoring", // RootPath specifies the root for asynqmon app
	// 	RedisConnOpt: asynq.RedisClientOpt{Addr: ":6379"},
	// })

	Init()

	logger.Logger.Info("setting up routes")

	Router = chi.NewRouter()

	Router.Use(middleware.Logger)

	// r.Handle(h.RootPath()+"/", h)

	Router.Get("/", func(w http.ResponseWriter, r *http.Request) {

		util.RespondWithJSON(w, http.StatusOK, "Hello, world")
	})

	currentDirectory, _ := os.Getwd()
	pathToDocFile := filepath.Join(currentDirectory, "docs")

	fs := http.FileServer(http.Dir(pathToDocFile)) // https://ribice.medium.com/serve-swaggerui-within-your-golang-application-5486748a5ed4
	Router.Handle("/docs/*", http.StripPrefix("/docs/", fs))

	routes.SetupRoutes(Router)

	var addr string

	if config.EnvConfig.Port == "" {
		addr = ":5000"
	} else {
		addr = ":" + config.EnvConfig.Port
	}

	logger.Logger.Info("server is running on port " + addr)
	http.ListenAndServe(addr, Router)
}
