package main

import (
	"db"
	"handlers"

	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// If it's Dev environment, load the .env file
	prod := os.Getenv("PROD")

	if prod != "true" {
		// load .env file
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	dbType := os.Getenv("DB_TYPE")

	db.SetDataBaseConnector(dbType)

	if !db.DbConn.IsConnection() {
		log.Fatal("No connection to the DB")
	}

	log.Println("Connection successful to the DB")

	handlers.Handlers()
}
