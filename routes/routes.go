package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sanyogpatel-tecblic/RBCA/endpoints"
	"github.com/sanyogpatel-tecblic/RBCA/middleware"
	"github.com/sanyogpatel-tecblic/RBCA/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Routes() {
	dsn := "host=localhost user=postgres dbname=student password=root port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	// Auto-migrate the User model to create the "users" table if it doesn't exist
	db.AutoMigrate(&models.User{})
	router := gin.Default()

	// // Apply RBAC to the endpoints
	// router.POST("/users", LoginHandler(db, []string{"admin"}), CreateUserHandler(db))
	router.DELETE("/users/:id", endpoints.DeleteUsers(db))
	router.GET("/users", middleware.AuthMiddleware(db), endpoints.GetUserHandler(db))
	router.POST("/users", middleware.AuthMiddleware(db), endpoints.CreateUserHandler(db))
	router.POST("/login", endpoints.LoginHandler(db), endpoints.CreateAccessTokenHandler)
	router.POST("/register", endpoints.Register(db))

	//admin
	router.GET("/users/requests", middleware.AuthMiddleware(db), endpoints.GetAllRequests(db))

	// Run the server
	router.Run(":8080")
}
