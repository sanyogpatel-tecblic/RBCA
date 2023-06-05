package main

import (
	"fmt"
	"log"

	"github.com/sanyogpatel-tecblic/RBCA/models"
	"github.com/sanyogpatel-tecblic/RBCA/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Define the User model

func main() {
	dsn := "host=localhost user=postgres dbname=student password=root port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is getting started...")
	fmt.Println("Listening at port 8080 ...")
	db.AutoMigrate(&models.User{})
	routes.Routes()
}
