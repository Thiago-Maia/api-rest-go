package main

import (
	"log"

	"github.com/Thiago-Maia/gin-api-rest-alura/database"
	"github.com/Thiago-Maia/gin-api-rest-alura/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}
	database.Connect()
	routes.HandleRequests()
}
