package main

import (
	"github.com/Adedunmol/mycart/internal/app"
	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/Adedunmol/mycart/internal/redis"
	"github.com/Adedunmol/mycart/internal/tasks"
)

const redisAddress = "127.0.0.1:6379"

func main() {
	logger.Logger.Info("app is running")

	go tasks.Init(redisAddress)
	defer tasks.Close()

	go redis.Init(redisAddress)
	defer redis.Close()

	logger.Logger.Info("app is running")
	app.Run()
}
