package main

import (
	"proyecto-integrador/app"
	"proyecto-integrador/db"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando .env: %v", err)
	}

	db.StartDbEngine()
	app.StartRoute()
}
