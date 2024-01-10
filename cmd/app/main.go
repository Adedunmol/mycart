package main

import (
	"log"
	"log/slog"

	"github.com/Adedunmol/mycart/internal/app"
	"github.com/Adedunmol/mycart/internal/database"
	"github.com/go-chi/httplog/v2"
	"github.com/joho/godotenv"
)

func Initializers() database.DbInstance {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.InitDB()

	if err != nil {
		log.Panic(err)
	}

	return db
}

func main() {
	logger := httplog.NewLogger("mycart-logs", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		MessageFieldName: "message",
	})

	app.Run(logger)
	// r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })
	// http.ListenAndServe(":3000", r)
}
