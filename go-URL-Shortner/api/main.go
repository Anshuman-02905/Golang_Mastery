package main

import (
	"fmt"
	"go-url-shortener/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error At go dot env")
		fmt.Println(err)
	}
	app := fiber.New()
	routes.SetupRoutes(app)
	app.Use(logger.New())
	log.Println("THIS IS APPPORT", os.Getenv("APP_PORT"))
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
