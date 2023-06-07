package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanyogpatel-tecblic/RBCA/email"
	"github.com/sanyogpatel-tecblic/RBCA/models"
	"gorm.io/gorm"
)

func GetUserHandler(db *gorm.DB) gin.HandlerFunc {
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

		var user models.User
		err = db.Where("id = ?", userID).First(&user).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		if user.Approved != 1 && user.Role != "admin" {
			// Send email alert
			email.SendEmailAlert(user.Email, "GetUserHandler API was called.....", "You do not have approval to access this feature. Please contact ADMINISTRATION for further process")

			c.JSON(http.StatusOK, gin.H{"message": "You do not have approval to access this feature. Please contact ADMINISTRATION for further process."})
			return
		}
		if user.Role != "admin" && user.Role != "manager" {
			email.SendEmailAlert(user.Email, "GetUserHandler API was called.....", "You are not allowed to access this feature. Please contact ADMINISTRATION for further asistance.")
			c.JSON(http.StatusOK, gin.H{"message": "You are not allowed to access this feature. Please contact ADMINISTRATION for further process."})
			return
		}

		// Retrieve the list of users from the database
		var users []models.User
		result := db.Find(&users)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}
