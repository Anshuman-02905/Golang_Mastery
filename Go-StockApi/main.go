package main

import (
	"fmt"
	"log"
	"net/http"
	"stock/config"
	"stock/routes"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize router with defined routes
	r := routes.SetupRoutes()
	// Server start message
	fmt.Println("🚀 Starting server on the port 8080...")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("❌ Failed to start server : %v", err)
	}

}
