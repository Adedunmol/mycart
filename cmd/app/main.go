package main

import (
	"log/slog"

	"github.com/Adedunmol/mycart/internal/app"
	"github.com/Adedunmol/mycart/internal/redis"
	"github.com/Adedunmol/mycart/internal/tasks"
	"github.com/go-chi/httplog/v2"
)

const redisAddress = "127.0.0.1:6379"

func main() {
	tasks.Init(redisAddress)
	defer tasks.Close()

	redis.Init(redisAddress)
	defer redis.Close()

	logger := httplog.NewLogger("mycart-logs", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		MessageFieldName: "message",
	})

	app.Run(logger)
}
