package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(dsn string) {
	/*
		Takes dsn Input and establishes Connection
	*/
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Connection NOT Established", err)
		return
	}

	DB = db
	fmt.Println("Connection Established")

}
