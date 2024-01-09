package main

import (
	"log/slog"

	"github.com/Adedunmol/mycart/internal/app"
	"github.com/go-chi/httplog/v2"
)

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
