package main

import (
	"log"

	"github.com/Adedunmol/mycart/internal/app"
	"github.com/Adedunmol/mycart/internal/config"
	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/Adedunmol/mycart/internal/redis"
	"github.com/Adedunmol/mycart/internal/tasks"
)

func main() {
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	go tasks.Init(config.EnvConfig.RedisAddress)

	go tasks.Run()

	defer tasks.Close()

	go redis.Init(config.EnvConfig.RedisAddress)
	defer redis.Close()

	logger.Logger.Info("app is running")
	app.Run()
}
