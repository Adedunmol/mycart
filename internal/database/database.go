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

func InitDB(config util.Config) (DbInstance, error) {
	var err error

	db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to the db: ", err)
	}
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running migrations")
	db.AutoMigrate(&models.Role{}, &models.User{})

	Database = DbInstance{DB: db}

	return Database, nil
}
