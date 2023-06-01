package main

import (
	"fmt"

	"github.com/sanyogpatel-tecblic/RBCA/routes"
)

// Define the User model

func main() {
	fmt.Println("Server is getting started...")
	fmt.Println("Listening at port 8080 ...")

	routes.Routes()
}
