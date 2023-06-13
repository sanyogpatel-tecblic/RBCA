package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/sanyogpatel-tecblic/RBCA/models"
	"gorm.io/gorm"
)

func GetUserHandler(db *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fetch the access token from the request header
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access token"})
			return
		}

		// Verify the access token
		claims, err := VerifyAccessToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Get the user ID from the claims
		userID, ok := claims["userID"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in access token"})
			return
		}

		// Check if the user data exists in Redis cache
		// ctx := context.Background()
		userData, err := redisClient.Get(userID).Result()
		if err == nil {
			// Cache hit: User data found in Redis
			log.Println("Data retrieved from cache")

			// Parse the cached data and return the response
			var users []models.User
			err = json.Unmarshal([]byte(userData), &users)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to parse cached data"})
				return
			}
			c.JSON(http.StatusOK, users)
			return
		} else if err != redis.Nil {
			// Error occurred while accessing Redis
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access Redis cache"})
			return
		}

		// Cache miss: User data not found in Redis
		log.Println("Data retrieved from the database")

		// Query the database to retrieve the user data
		var users []models.User
		result := db.Find(&users)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		// Cache the retrieved user data in Redis
		usersData := convertUsersToJSON(users)
		err = redisClient.Set(userID, usersData, time.Hour).Err()
		if err != nil {
			log.Println("Failed to cache user data:", err)
		}

		c.JSON(http.StatusOK, users)
	}
}

func convertUsersToJSON(users []models.User) string {
	// Convert users slice to JSON
	jsonData, err := json.Marshal(users)
	if err != nil {
		log.Println("Failed to convert users to JSON:", err)
		return ""
	}

	// Return the JSON representation of users
	return string(jsonData)
}
