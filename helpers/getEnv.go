package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

/* GetEnvVariable Load the .env file and returns the value of the key */
func GetEnvVariable(key string) string {

	// If it's Dev environment, load the .env file
	prod := os.Getenv("PROD")

	if prod != "true" {
		// load .env file
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return os.Getenv(key)
}
