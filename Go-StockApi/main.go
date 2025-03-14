package main

import (
	"fmt"
	"log"
	"net/http"
	"stock/config"
	"stock/routes"
)

func main() {

	config.LoadEnv()
	r := routes.SetupRoutes()
	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))

}
