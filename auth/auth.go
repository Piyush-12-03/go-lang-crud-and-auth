package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/go-project/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var jwtSecret = []byte("KJKJvjVJgj&^574&768&*^&$728y7JvjVJFjgvjhgVuyglwajhqoiewosiqwhaiVUVKUVKJhw")
var userService *services.UsersService

// GenerateJWT generates a JWT token for a user including role
func GenerateJWT(userId int, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"role":    role,                                  // Add role to the token
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiry time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates the JWT token
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

const (
	key    = "uygs@*ibiIVSUYU@sIUSbibspougefuvASDSUGU@*W&873ni3h993oBIsib2"
	MaxAge = 86400 * 30 // Session expiration
	IsProd = false      // Whether it's production or not
)

func NewAuth(router *gin.Engine, userService *services.UsersService) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Get Google credentials from environment variables
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	// Set up session store
	store := sessions.NewCookieStore([]byte(key))
	gothic.Store = store

	// Register Google provider
	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:8888/auth/callback/google", "email", "profile"),
	)

	// Authentication routes
	router.GET("/auth/google", func(c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", "google"))
		gothic.BeginAuthHandler(c.Writer, c.Request)
	})

	router.GET("/auth/callback/google", func(c *gin.Context) {
		getAuthCallBackFunctions(c, userService)
	})
}

// Callback function
func getAuthCallBackFunctions(c *gin.Context, userService *services.UsersService) {
	// Get the provider name
	provider, err := gothic.GetProviderName(c.Request)
	if err != nil {
		fmt.Println("Error getting provider name:", err, provider)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get provider"})
		return
	}

	// Complete user authentication
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Println("Error completing user authentication:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to complete authentication"})
		return
	}

	// Extract email from the Google user payload
	email := user.Email

	fmt.Println("Google Email : ", email)
	// Check if the email exists in the user repository
	existingUser, err := userService.FindUserByEmail(email)
	if err != nil {
		fmt.Println("Error searching for user in repository:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching user"})
		return
	}

	if existingUser != nil {
		// User exists, generate a JWT token
		token, err := GenerateJWT(existingUser.Id, existingUser.Email, existingUser.Role)
		if err != nil {
			fmt.Println("Error generating JWT:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		// Return the user and token
		c.JSON(http.StatusOK, gin.H{
			"user":  existingUser,
			"token": token,
		})
	} else {
		// User does not exist, prompt for registration
		c.JSON(http.StatusOK, gin.H{
			"message": "You need to register to this platform",
		})
	}
}
