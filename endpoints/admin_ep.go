package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanyogpatel-tecblic/RBCA/models"
	"gorm.io/gorm"
)

func Approve(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

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
