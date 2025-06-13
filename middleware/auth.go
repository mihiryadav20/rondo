package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	
	"rondo/utils"
)

// AuthMiddleware is a middleware for JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		
		// Check if header is empty
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		
		// Check if header format is valid
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}
		
		// Extract token
		tokenString := parts[1]
		
		// Validate token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		
		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("phone", claims.Phone)
		c.Set("firstName", claims.FirstName)
		c.Set("lastName", claims.LastName)
		
		c.Next()
	}
}
