package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := "host=db user=postgres password=postgres dbname=go_demo port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Database connection successful")
	return nil
}
