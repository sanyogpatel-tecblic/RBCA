package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanyogpatel-tecblic/RBCA/email"
	"github.com/sanyogpatel-tecblic/RBCA/models"
	"gorm.io/gorm"
)

func AddReactTopics(db *gorm.DB) gin.HandlerFunc {
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
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in access token"})
			return
		}

		var newtopic models.Topic

		if err := c.ShouldBindJSON(&newtopic); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if newtopic.Topic == "" {
			apierror := models.APIError{
				Code:    http.StatusBadRequest,
				Message: "Topic is required",
			}
			c.Header(apierror.Message, "")
			json.NewEncoder(c.Writer).Encode(apierror)
			return
		}
		var user models.User
		err = db.Where("id = ?", userID).First(&user).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		if user.Role != "admin" && user.Role != "react" {
			email.SendEmailAlert(user.Email, "AddReactTopics API was Called....", "You are not allowed to access this feature. Please contact ADMINISTRATION for further asistance.")
			c.JSON(http.StatusOK, gin.H{"message": "You are not allowed to access this feature. Please contact ADMINISTRATION for further process."})
			return
		}

		if user.Approved != 1 {
			// Send email alert
			email.SendEmailAlert(user.Email, "AddReactTopics API was Called....", "You do not have approval to access this feature. Please contact ADMINISTRATION for further process")

			c.JSON(http.StatusOK, gin.H{"message": "You do not have approval to access this feature. Please contact ADMINISTRATION for further process."})
			return
		}

		newtopic.UserID = userIDInt

		result := db.Table("topics").Create(&newtopic)
		if result.Error != nil {
			apierror := models.APIError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create user: " + result.Error.Error(),
			}
			c.Header(apierror.Message, "")
			json.NewEncoder(c.Writer).Encode(apierror)
			return
		}
		cacheKey := "reactTopics"
		delete(Cache, cacheKey)

		// Return success response
		c.JSON(http.StatusOK, newtopic)
	}
}

var Cache map[string]interface{}

func init() {
	Cache = make(map[string]interface{})
}
func GetReactTopicsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		cacheKey := "reactTopics"

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

		err = db.Where("id = ?", userID).First(&user).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		if user.Approved != 1 && user.Role != "react" {
			// Send email alert
			email.SendEmailAlert(user.Email, "GetReactTopicsHandler API was Called....", "You do not have approval to access this feature. Please contact ADMINISTRATION for further process")

			c.JSON(http.StatusOK, gin.H{"message": "You do not have approval to access this feature. Please contact ADMINISTRATION for further process."})
			return
		}

		if user.Role != "admin" && user.Role != "react" {
			email.SendEmailAlert(user.Email, "GetReactTopicsHandler API was Called....", "You are not allowed to access this feature. Please contact ADMINISTRATION for further assistance.")
			c.JSON(http.StatusOK, gin.H{"message": "You are not allowed to access this feature. Please contact ADMINISTRATION for further process."})
			return
		}

		// Check if data exists in the cache
		if cachedData, ok := Cache[cacheKey]; ok {
			c.JSON(http.StatusOK, cachedData)
			log.Println("from cache......................................................")
			email.SendEmailAlert(user.Email, "GetReactTopicsHandler API was Called....", "")
			return
		}

		// Retrieve the list of react topics from the database
		var topics []models.Topic
		result := db.Joins("JOIN users ON users.id = topics.user_id").Where("users.role = ?", "react").Find(&topics)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve topics"})
			return
		}
		Cache[cacheKey] = topics

		cacheExpiration := time.Now().Add(5 * time.Minute)
		Cache[cacheKey+"Expiration"] = cacheExpiration
		fmt.Println("from database,,,......................................................")
		email.SendEmailAlert(user.Email, "GetReactTopicsHandler API was Called....", "")

		c.JSON(http.StatusOK, topics)
	}
}

func GetNodeTopicsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fetch the access token from the request header
		var user models.User
		cacheKey := "nodeTopics"

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

		err = db.Where("id = ?", userID).First(&user).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		if user.Approved != 1 && user.Role != "admin" {
			// Send email alert
			email.SendEmailAlert(user.Email, "GetNodeTopicsHandler API was Called....", "You do not have approval to access this feature. Please contact ADMINISTRATION for further process")
			c.JSON(http.StatusOK, gin.H{"message": "You do not have approval to access this feature. Please contact ADMINISTRATION for further process."})
			return
		}

		if user.Role != "admin" && user.Role != "nodejs" {
			email.SendEmailAlert(user.Email, "GetNodeTopicsHandler API was Called....", "You are not allowed to access this feature. Please contact ADMINISTRATION for further assistance.")
			c.JSON(http.StatusOK, gin.H{"message": "You are not allowed to access this feature. Please contact ADMINISTRATION for further process."})
			return
		}
		if cachedData, ok := Cache[cacheKey]; ok {
			c.JSON(http.StatusOK, cachedData)
			log.Println("from cache,,,......................................................")
			email.SendEmailAlert(user.Email, "GetNodeTopicsHandler API was Called....", "")
			return
		}

		// Retrieve the list of react topics from the database
		var topics []models.Topic
		result := db.Joins("JOIN users ON users.id = topics.user_id").Where("users.role = ?", "nodejs").Find(&topics)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve topics"})
			return
		}
		Cache[cacheKey] = topics

		cacheExpiration := time.Now().Add(5 * time.Minute)
		Cache[cacheKey+"Expiration"] = cacheExpiration
		fmt.Println("from database,,,......................................................")

		email.SendEmailAlert(user.Email, "GetNodeTopics API was Called....", "")

		c.JSON(http.StatusOK, topics)
	}
}

func AddNodeTopics(db *gorm.DB) gin.HandlerFunc {
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
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in access token"})
			return
		}

		var newtopic models.Topic

		if err := c.ShouldBindJSON(&newtopic); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if newtopic.Topic == "" {
			apierror := models.APIError{
				Code:    http.StatusBadRequest,
				Message: "Topic is required",
			}
			c.Header(apierror.Message, "")
			json.NewEncoder(c.Writer).Encode(apierror)
			return
		}
		var user models.User
		err = db.Where("id = ?", userID).First(&user).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		if user.Role != "admin" && user.Role != "nodejs" {
			email.SendEmailAlert(user.Email, "AddNodeTopics API was Called....", "You are not allowed to access this feature. Please contact ADMINISTRATION for further asistance.")
			c.JSON(http.StatusOK, gin.H{"message": "You are not allowed to access this feature. Please contact ADMINISTRATION for further process."})
			return
		}

		if user.Approved != 1 {
			// Send email alert
			email.SendEmailAlert(user.Email, "AddNodeTopics API was Called....", "You do not have approval to access this feature. Please contact ADMINISTRATION for further process")

			c.JSON(http.StatusOK, gin.H{"message": "You do not have approval to access this feature. Please contact ADMINISTRATION for further process."})
			return
		}

		newtopic.UserID = userIDInt

		result := db.Table("topics").Create(&newtopic)
		if result.Error != nil {
			apierror := models.APIError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create user: " + result.Error.Error(),
			}
			c.Header(apierror.Message, "")
			json.NewEncoder(c.Writer).Encode(apierror)
			return
		}

		// Return success response
		c.JSON(http.StatusOK, newtopic)
	}
}

func CheckCacheHandler(c *gin.Context) {
	cacheKeys := []string{"reactTopics", "nodeTopics"}
	cachedData := make(map[string]interface{})

	for _, cacheKey := range cacheKeys {
		if data, ok := Cache[cacheKey]; ok {
			cachedData[cacheKey] = data
		}
	}

	if len(cachedData) > 0 {
		c.JSON(http.StatusOK, cachedData)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Data not found in cache"})
}
