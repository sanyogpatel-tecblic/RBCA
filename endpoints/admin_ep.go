package endpoints

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanyogpatel-tecblic/RBCA/email"
	"github.com/sanyogpatel-tecblic/RBCA/models"
	"gorm.io/gorm"
)

func Approve(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var users models.User
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

		result := db.Where("id = ?", userID).Find(&users)

		if users.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		var user models.User

		id := c.Param("id")

		result = db.First(&user, id)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		user.Approved = 1

		result = db.Save(&user)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		email.SendEmailAlert(user.Email, "You Are Approved! ", "Now You Can Access APIs")
		c.JSON(http.StatusOK, gin.H{"message": "User Approoved Successfully"})
	}
}
func Reject(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var users models.User
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

		result := db.Where("id = ?", userID).Find(&users)

		if users.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		var user models.User

		id := c.Param("id")

		result = db.First(&user, id)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		user.Approved = 1

		result = db.Save(&user)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		email.SendEmailAlert(user.Email, "You Are Approved! ", "Now You Can Access APIs")
		c.JSON(http.StatusOK, gin.H{"message": "User Approoved Successfully"})
	}
}

func GetAllRequests(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		var user models.User
		err = db.Where("id = ?", userID).First(&user).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		// Retrieve the list of users where approved is 0
		var users []models.User
		result := db.Where("approved = ?", 0).Find(&users)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		// Print user IDs to console
		for _, user := range users {
			fmt.Println("User ID:", user.ID)
		}

		c.JSON(http.StatusOK, users)
	}
}

func CreateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authenticated user from the context
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

		var user2 models.User
		err = db.Where("id = ?", userID).First(&user2).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		// Parse the request body to get the user data
		var newUser models.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if user2.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		// Create the user in the database
		result := db.Create(&newUser)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, newUser)
	}
}
