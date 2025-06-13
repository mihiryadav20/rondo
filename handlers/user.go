package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	
	"rondo/models"
	"rondo/utils"
)

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var req models.UserRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	
	// Validate date format
	_, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	
	// Get phone number from JWT token
	phone, exists := c.Get("phone")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number not found in token"})
		return
	}
	
	phoneNumber, ok := phone.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid phone number format in token"})
		return
	}
	
	// Check if user already exists
	_, userExists := utils.GetUserByPhone(phoneNumber)
	if userExists {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this phone number already exists"})
		return
	}
	
	// Create user
	user, err := utils.CreateUser(req, phoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return user data with token
	c.JSON(http.StatusCreated, gin.H{
		"user": models.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			DOB:       user.DOB,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
		},
		"token": token,
	})
}

// GetUserProfile retrieves user profile by phone number
func GetUserProfile(c *gin.Context) {
	phone := c.Param("phone")
	
	user, exists := utils.GetUserByPhone(phone)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	c.JSON(http.StatusOK, models.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		DOB:       user.DOB,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	})
}
