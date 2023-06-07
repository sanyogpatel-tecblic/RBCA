package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanyogpatel-tecblic/RBCA/email"
	"github.com/sanyogpatel-tecblic/RBCA/models"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newuser models.User

		if err := c.ShouldBindJSON(&newuser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if newuser.Username == "" {
			apierror := models.APIError{
				Code:    http.StatusBadRequest,
				Message: "username is required",
			}
			c.Header(apierror.Message, "")
			json.NewEncoder(c.Writer).Encode(apierror)
			return
		}
		if newuser.Password == "" {
			apierror := models.APIError{
				Code:    http.StatusBadRequest,
				Message: "password is required",
			}
			c.Header(apierror.Message, "password is required 2")
			json.NewEncoder(c.Writer).Encode(apierror)
			return
		}

		// Set the approved field to 0 (false) by default
		newuser.Approved = 0

		result := db.Table("users").Create(&newuser)
		if result.Error != nil {
			apierror := models.APIError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create user: " + result.Error.Error(),
			}
			c.Header(apierror.Message, "")
			json.NewEncoder(c.Writer).Encode(apierror)
			return
		}
		email.SendEmailAlert2("ad2491min@gmail.com", "ad2491min@gmail.com", "New user have just registered please have a look at your pending Requests!", newuser.Username)

		email.SendEmailAlert(newuser.Email, "Register API was Called....", "You have succesfully registered, please wait till admin approves you!")

		// Return success response
		c.JSON(http.StatusOK, newuser)

	}
}
