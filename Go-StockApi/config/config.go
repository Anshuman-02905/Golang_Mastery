package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from the .env file located in the "config" directory.
// It logs a fatal error and stops execution if the file cannot be loaded.
func LoadEnv() {

	// Attempt to load the .env file
	if err := godotenv.Load("config/.env"); err != nil {
		log.Fatal("❌ Error loading .env file")
	}
}

/*
GetEnv retrives the value from enviroment Variables
Args:

	-key (string): The name of the enviroment variable to retrieve

Returns:
  - string: The value of the mentioned enviroment variable or an empty string
*/
func GetEnv(key string) string {
	return os.Getenv(key)
}
