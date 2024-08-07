package main

import (
	"log"

	"github.com/Adedunmol/mycart/internal/app"
	"github.com/Adedunmol/mycart/internal/config"
	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/Adedunmol/mycart/internal/redis"
	"github.com/Adedunmol/mycart/internal/tasks"
)

const redisAddress = "127.0.0.1:6379"

func main() {
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	go tasks.Init(redisAddress)

	go tasks.Run()

	defer tasks.Close()

	go redis.Init(redisAddress)
	defer redis.Close()

	logger.Logger.Info("app is running")
	app.Run()
}
