package database

import (
	"fmt"
	"log"
	"os"

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

	fmt.Println(os.Getenv("DATABASE_URL"))
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to the db: ", err)
	}
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running migrations")
	db.AutoMigrate()

	Database = DbInstance{DB: db}

	return Database, nil
}
