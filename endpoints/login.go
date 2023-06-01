package endpoints

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sanyogpatel-tecblic/RBCA/models"
	"gorm.io/gorm"
)

func GenerateAccessToken(userID string) (string, error) {
	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID

	// Generate encoded token and return it
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyAccessToken(accessToken string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		// Return the secret key
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Get the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}

func GetUserIDFromAccessToken(tokenString string) (string, error) {
	// Parse the token with your JWT secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Replace "your-secret-key" with your actual secret key
		return []byte("your-secret-key"), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to extract user ID from token")
	}
	userID, ok := claims["userID"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract user ID from token")
	}

	return userID, nil
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var reqBody struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Query the database for the user
		var user models.User
		result := db.Table("users").Where("username = ? AND password = ?", reqBody.Username, reqBody.Password).First(&user)
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		} else if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		// Store the authenticated user in the context for future use
		c.Set("user", user)

		// Continue to the next middleware or handler
		c.Next()
	}
}

func CreateAccessTokenHandler(c *gin.Context) {
	// Get the authenticated user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve authenticated user"})
		return
	}

	// Generate an access token for the user
	accessToken, err := GenerateAccessToken(strconv.Itoa(int(user.(models.User).ID)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send the access token in the response
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
