package helpers

import (
	"os"
)

/* GetEnvVariable Load the .env file and returns the value of the key */
func GetEnvVariable(key string) string {
	// load .env file
	// err := godotenv.Load()

	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	return os.Getenv(key)
}
