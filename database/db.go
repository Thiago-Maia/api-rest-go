package database

import (
	"log"
	"os"

	"github.com/Thiago-Maia/gin-api-rest-alura/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Db  *gorm.DB
	err error
)

func Connect() {

	dsn := os.Getenv("CONNECTION_STRING")
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar com banco de dados")
	}
	Db.AutoMigrate(&models.Student{})
}
