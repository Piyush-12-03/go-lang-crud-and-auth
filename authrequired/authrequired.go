package authrequired

import (
	"net/http"
	"strings"

	"example.com/go-project/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// RoleBasedAuth checks for JWT token and verifies user roles
func RoleBasedAuth(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Token Required."})
			c.Abort()
			return
		}

		// Check if the token is prefixed with "Bearer"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// Validate the token
		token, err := auth.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract user ID and role from token claims
			userId := claims["user_id"]
			role := claims["role"].(string)

			// Set user ID in context
			c.Set("user_id", userId)

			// Check if user has the required role
			if role != requiredRole {
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient privileges"})
				c.Abort()
				return
			}

			// Continue to the next handler
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}
