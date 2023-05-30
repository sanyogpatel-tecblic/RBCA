package main

import "log"

func main() {

	dsn := "host=localhost user=postgres dbname=student password=root port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}
