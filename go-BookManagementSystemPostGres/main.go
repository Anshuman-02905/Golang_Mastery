package main

import (
	"bookstore/config"
	"bookstore/routes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Book struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"username" gorm:"column:username"`
	AuthorName string    `json:"author" gorm:"column:author"`
	Genre      string    `json:"genre" gorm:"column:genre"`
	Rent       string    `json:"rent" gorm:"column:rent"`
	Currency   string    `json:"currency" gorm:"column:currency"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"` // Auto-set on insert
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"` // Auto-set on update
}

func main() {

	content, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		fmt.Println("Error when opening file", err)
	}
	var result map[string]interface{}

	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", result["HOST"], result["POSTGRES_USER"],
		result["POSTGRES_PASSWORD"], result["POSTGRES_DB"], result["PORT"], result["SSL_MODE"])

	config.ConnectDatabase(dsn)

	r := routes.SetupRoutes()
	r.Run(":8080")

}
