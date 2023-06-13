package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/sanyogpatel-tecblic/RBCA/endpoints"
	"github.com/sanyogpatel-tecblic/RBCA/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Routes() {
	dsn := "host=localhost user=postgres dbname=student password=root port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	// Apply RBAC to the endpoints
	// router.POST("/users", LoginHandler(db, []string{"admin"}), CreateUserHandler(db))
	router.DELETE("/users/:id", endpoints.DeleteUsers(db))
	router.GET("/users", middleware.AuthMiddleware(db), endpoints.GetUserHandler(db, redisClient))
	router.POST("/users", middleware.AuthMiddleware(db), endpoints.CreateUserHandler(db))
	router.POST("/login", endpoints.LoginHandler(db), endpoints.CreateAccessTokenHandler)
	router.POST("/register", endpoints.Register(db))

	//admin
	router.GET("/users/requests", middleware.AuthMiddleware(db), endpoints.GetAllRequests(db))
	router.PATCH("/users/approve/:id", middleware.AuthMiddleware(db), endpoints.Approve(db))

	//react
	router.POST("/topic/react", middleware.AuthMiddleware(db), endpoints.AddReactTopics(db))
	router.GET("/topic/react", middleware.AuthMiddleware(db), endpoints.GetReactTopicsHandler(db))

	//node
	router.POST("/topic/node", middleware.AuthMiddleware(db), endpoints.AddNodeTopics(db))
	router.GET("/topic/node", middleware.AuthMiddleware(db), endpoints.GetNodeTopicsHandler(db))

	router.GET("/check-cache", endpoints.CheckCacheHandler)

	// Run the server
	router.Run(":8080")
}
