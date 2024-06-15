package main

import (
	"log/slog"

	"github.com/Adedunmol/mycart/internal/app"
	"github.com/go-chi/httplog/v2"
	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	client := asynq.NewClient((asynq.RedisClientOpt{Addr: redisAddr}))
	defer client.Close()

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
