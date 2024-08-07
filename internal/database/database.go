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

func InitDB(databaseUrl string) error {
	var err error

	DB, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{TranslateError: true, Logger: logger.Default.LogMode(logger.Silent)})

	if err != nil {
		customLogger.Logger.Error("error connecting to db: ")
		customLogger.Logger.Error(err.Error())

		return err
	}

	if config.EnvConfig.Environment != "test" {
		// DB.Logger = logger.Default.LogMode(logger.Info)

		customLogger.Logger.Info("Running migrations")
	}

	DB.AutoMigrate(&models.Otp{})
	DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Review{}, &models.Product{}, &models.Order{}, &models.CartItem{}, &models.Cart{})

	return nil
}
