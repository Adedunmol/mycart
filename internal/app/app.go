package app

import (
	"log"

	"github.com/Adedunmol/mycart/internal/database"
)

func Initializers() database.DbInstance {
	db, err := database.InitDB()

	if err != nil {
		log.Panic(err)
	}

	return db
}

func Run() {

	Initializers()
}
