package main

import (
	"go-BookManagementSystem/config"
	"go-BookManagementSystem/models"
	"go-BookManagementSystem/routes"
)

func main() {
	config.ConnectDatabase("database/Library.db")

	config.DB.AutoMigrate(&models.Book{})
	r := routes.SetupRoutes()
	r.Run(":8080")
}
