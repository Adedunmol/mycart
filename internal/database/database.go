package database

import (
	"log"

	"github.com/Adedunmol/mycart/internal/config"
	"github.com/Adedunmol/mycart/internal/models"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	DB *gorm.DB
}

var Database DbInstance

func InitDB() (DbInstance, error) {
	var err error
	var db *gorm.DB

	if config.EnvConfig.Environment != "test" {
		db, err = gorm.Open(postgres.Open(config.EnvConfig.DatabaseUrl), &gorm.Config{})
	}

	db, err = gorm.Open(postgres.Open(config.EnvConfig.TestDatabaseUrl), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to the db: ", err)
	}

	if config.EnvConfig.Environment != "test" {
		db.Logger = logger.Default.LogMode(logger.Info)

		log.Println("Running migrations")
	}

	db.AutoMigrate(&models.Role{}, &models.User{})

	Database = DbInstance{DB: db}

	return Database, nil
}
