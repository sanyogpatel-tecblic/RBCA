package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanyogpatel-tecblic/RBCA/models"
	"gorm.io/gorm"
)

func DeleteUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var user models.User
		result := db.First(&user, id)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		user.Deleted = 1
		result = db.Save(&user)

		result = db.Delete(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User Deleted Successfully"})
	}
}
