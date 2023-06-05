package endpoints

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanyogpatel-tecblic/RBCA/models"
	"gopkg.in/gomail.v2"
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

		if user.Approved == 0 && user.Role != "admin" && user.Role != "manager" {
			// Send email alert
			sendEmailAlert(user.Email, "You do not have approval to access this feature. Please contact ADMINISTRATION for further process")

			c.JSON(http.StatusOK, gin.H{"message": "You do not have approval to access this feature. Please contact ADMINISTRATION for further process."})
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

func sendEmailAlert(email string, error string) {
	// Set up SMTP connection details for Gmail
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	senderEmail := "sanyogp.249@gmail.com"
	senderPassword := "udkaoitkaqgfipyh"
	recipientEmail := "21msit037@charusat.edu.in"

	// Create a new email message
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "API Alert")

	// Get the current time with microseconds
	timestamp := time.Now().Format("2006-01-02 15:04:05.000000")

	// Add the timestamp to the email body
	body := "The GetUserHandler API was called at " + timestamp + error
	m.SetBody("text/plain", body)

	// Set up the SMTP authentication
	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		// Handle the error, e.g., log it or return an error response
		log.Println("Failed to send email alert:", err)
	}
}
