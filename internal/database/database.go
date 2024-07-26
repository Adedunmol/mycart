package database

import (
	"github.com/Adedunmol/mycart/internal/config"
	customLogger "github.com/Adedunmol/mycart/internal/logger"
	"github.com/Adedunmol/mycart/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error

	if config.EnvConfig.Environment == "test" {
		DB, err = gorm.Open(postgres.Open(config.EnvConfig.TestDatabaseUrl), &gorm.Config{TranslateError: true})
	} else {
		DB, err = gorm.Open(postgres.Open(config.EnvConfig.DatabaseUrl), &gorm.Config{TranslateError: true})
	}

	if err != nil {
		customLogger.Logger.Error("error connecting to db: ")
		customLogger.Logger.Error(err.Error())
	}

	if config.EnvConfig.Environment != "test" {
		DB.Logger = logger.Default.LogMode(logger.Info)

		customLogger.Logger.Info("Running migrations")
	}

	DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Review{}, &models.Product{}, &models.Order{}, &models.CartItem{}, &models.Cart{})
}
