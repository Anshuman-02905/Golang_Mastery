package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"username" gorm:"column:username"`
	AuthorName string    `json:"author" gorm:"column:author"`
	Genre      string    `json:"genre" gorm:"column:genre"`
	Rent       string    `json:"rent" gorm:"column:rent"`
	Currency   string    `json:"currency gorm:"column:currency"`
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Connection NOT Established")
	}

	err = db.Migrator().DropTable(&User{})

	db.AutoMigrate(&User{})

	user := User{
		Name:       "John Doe",
		AuthorName: "J.K. Rowling",
		Genre:      "Fiction",
		Rent:       "10",
		Currency:   "USD",
	}

	db.Create(&user)

	fmt.Println("Inserted USER ID", user.ID)
	fmt.Println("Created At:", user.CreatedAt)
	fmt.Println("Updated At:", user.UpdatedAt)

}
