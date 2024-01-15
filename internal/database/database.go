package database

import (
	"log"

	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/util"
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

	if util.EnvConfig.Environment != "test" {
		db, err = gorm.Open(postgres.Open(util.EnvConfig.DatabaseUrl), &gorm.Config{})
	}

	db, err = gorm.Open(postgres.Open(util.EnvConfig.TestDatabaseUrl), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to the db: ", err)
	}

	if util.EnvConfig.Environment != "test" {
		db.Logger = logger.Default.LogMode(logger.Info)

		log.Println("Running migrations")
	}

	db.AutoMigrate(&models.Role{}, &models.User{})

	Database = DbInstance{DB: db}

	return Database, nil
}
