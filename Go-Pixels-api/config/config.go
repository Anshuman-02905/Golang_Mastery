package config

import (
	"log"
	"net/http"
	"os"

	"github.com/Anshuman-02905/Golang-Mastery/Go-Pixels-api/models"

	"github.com/joho/godotenv"
)

func NewClient() *models.Client {
	c := http.Client{}
	token := GetKey("config/.env")
	return &models.Client{
		Token: token, HC: c,
	}
}

func GetKey(envpath string) string {
	err := godotenv.Load(envpath)
	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	s3cretkey := os.Getenv("Pexels_TOKEN")
	return s3cretkey
}
