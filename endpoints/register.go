package endpoints

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sanyogpatel-tecblic/RBCA/email"
	"github.com/sanyogpatel-tecblic/RBCA/models"
	"github.com/sfreiberg/gotwilio"
	"gorm.io/gorm"
)

const (
	twilioAccountSID  = "ACe1005348d84ee8b6152aaedeef0dfbd1"
	twilioAuthToken   = "d8a5aa894cd228d18b4dab33826061ad"
	twilioPhoneNumber = "+1 361 301 1373"
)

// Store OTP for each registration request
var otpStore = make(map[string]string)

func generateOTP() string {
	otp := rand.Intn(900000) + 100000
	return strconv.Itoa(otp)
}

func sendOTP(phoneNumber string, otp string) error {
	twilioClient := gotwilio.NewTwilioClient(twilioAccountSID, twilioAuthToken)
	_, _, err := twilioClient.SendSMS(twilioPhoneNumber, phoneNumber, "Your OTP for registration is: "+otp, "", "")

	return err
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newuser models.User

		if err := c.ShouldBindJSON(&newuser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if OTP is provided in the query parameters
		enteredOTP := c.Query("otp")
		if enteredOTP == "" {
			// Generate OTP
			otp := generateOTP()

			// Store the generated OTP for this registration request
			otpStore[newuser.Mobile] = otp

			// Send OTP to the user's mobile number
			err := sendOTP(newuser.Mobile, otp)

			if err != nil {
				apierror := models.APIError{
					Code:    http.StatusInternalServerError,
					Message: "Failed to send OTP: " + err.Error(),
				}
				c.JSON(http.StatusInternalServerError, apierror)
				return
			}
			email.SendEmailAlert(newuser.Email, "Register API was called ", "Your OTP for registration is: "+otp)
			c.JSON(http.StatusOK, gin.H{"message": "OTP sent. Please provide the OTP to complete registration."})
			return
		}

		// OTP provided, verify it
		expectedOTP, ok := otpStore[newuser.Mobile]
		if !ok || enteredOTP != expectedOTP {
			apierror := models.APIError{
				Code:    http.StatusBadRequest,
				Message: "Invalid OTP",
			}
			c.JSON(http.StatusBadRequest, apierror)
			return
		}

		newuser.Approved = 0

		result := db.Table("users").Create(&newuser)
		if result.Error != nil {
			apierror := models.APIError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create user: " + result.Error.Error(),
			}
			c.JSON(http.StatusInternalServerError, apierror)
			return
		}

		email.SendEmailAlert2("ad2491min@gmail.com", "ad2491min@gmail.com", "New user has just registered, please have a look at your pending Requests!", newuser.Username)
		email.SendEmailAlert(newuser.Email, "Register API was called", "You have successfully registered. Please wait until the admin approves you.")

		delete(otpStore, newuser.Mobile)
		c.JSON(http.StatusOK, newuser)
	}
}
