# GO CRUD Development with Docker and PostgreSQL

# How to Run the GO CRUD Application

## Step 1: Update the Environment File
Make sure your `config/.env` file is correctly set with the running container name and database name:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=root
DB_NAME=stock_db
```

## Step 2: Start the Docker Container
Run the following command to start the PostgreSQL container:

```
make dc-up
```

## Step 3: Create the Database Table
Run the following command to create the `stock` table:

```
make db-create
```

## Step 4: Verify Database Connection
You can check if the database is running and connect to it using:

```
make psql
```

## Step 5: Drop the Table (if needed)
To reset the database by dropping the `stock` table, run:

```
make db-drop
```

## Step 6: Install Dependencies
Ensure Go dependencies are installed and tidy up the module files:

```
go mod tidy
```

## Step 7: Run the Application
Finally, start the Go application:

```
go run main.go
```

Your CRUD API should now be running on `http://localhost:8080`! 🚀



## Introduction
This article walks through building a CRUD API in Go with Gorilla Mux, PostgreSQL, and Docker. We will cover environment setup, routing, middleware, database interactions, and Docker-based PostgreSQL integration.

---

## Project Structure
```
stock/
│-- config/
│   ├── .env
│   ├── config.go
│-- models/
│   ├── stock.go
│-- routes/
│   ├── routes.go
│-- middleware/
│   ├── handlers.go
│-- main.go
│-- docker-compose.yml
│-- Makefile
```

---

## Main File (main.go)
```go
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
```

---

## Configuration (config/config.go)
```go
package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Fatal("❌ Error loading .env file")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
```

---

## Models (models/stock.go)
```go
package models

type Stock struct {
	StockID int64  `json:"stockid"`
	Name    string `json:"name"`
	Price   string `json:"price"`
	Company string `json:"company"`
}
```

---

## Routes (routes/routes.go)
```go
package routes

import (
	"stock/middleware"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/stock/{id}", middleware.GetStock).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/stock/", middleware.GetAllStock).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/newstock", middleware.CreateStock).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/stock/{id}", middleware.UpdateStock).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/deletestock/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")

	return r
}
```

---

## Middleware (middleware/handlers.go)
```go
package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"stock/config"
	"stock/models"
	"strconv"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func connectDatabase() *sql.DB {
	DB_HOST := config.GetEnv("DB_HOST")
	DB_PORT := config.GetEnv("DB_PORT")
	DB_USER := config.GetEnv("DB_USER")
	DB_PASSWORD := config.GetEnv("DB_PASSWORD")
	DB_NAME := config.GetEnv("DB_NAME")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Database connection failed: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("❌ Database ping failed: ", err)
	}
	fmt.Println("✅ Successfully connected to the database")
	return db
}
```

---

## Docker Configuration (docker-compose.yml)
```yaml
version: '3.9'
services:
  postgres:
    image: postgres
    container_name: postgres_stock
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_PASSWORD=root
      - POSTGRES_USER=root
      - POSTGRES_DB=stock_db
```

---

## Database Migration (Makefile)
```make
DB_CONTAINER=postgres_stock
DB_NAME=stock_db
DB_USER=root

dc-up:
	docker compose up -d

dc-down:
	docker compose down

psql:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

db-create:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "CREATE TABLE IF NOT EXISTS stock (stockid BIGSERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, price VARCHAR(255) NOT NULL, company VARCHAR(255) NOT NULL);"

db-drop:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "DROP TABLE IF EXISTS stock;"
```

---

## Conclusion
This article covered setting up a CRUD API in Go using Gorilla Mux, PostgreSQL, and Docker. We structured the project, connected it to PostgreSQL, and managed migrations. This setup provides a scalable foundation for production-ready applications.

