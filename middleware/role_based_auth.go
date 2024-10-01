package middleware

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
			// Extract role from token claims
			role := claims["role"].(string)

			// Set user ID in context (optional if needed for further processing)
			c.Set("user_id", claims["user_id"])

			// Check if user has the required role
			if role != requiredRole {
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient privileges"})
				c.Abort()
				return
			}

			// Continue to the next handler if authorized
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}
