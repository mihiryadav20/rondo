package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	
	"rondo/models"
	"rondo/utils"
)

// TwilioClient is the global Twilio client
var TwilioClient *utils.TwilioClient

// InitHandlers initializes the handlers
func InitHandlers(twilioClient *utils.TwilioClient) {
	TwilioClient = twilioClient
}

// RequestOTP handles OTP request
func RequestOTP(c *gin.Context) {
	var req models.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	otp := utils.GenerateOTP()
	utils.StoreOTP(req.PhoneNumber, otp)

	// Send OTP via Twilio
	if err := TwilioClient.SendOTP(req.PhoneNumber, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// VerifyOTP handles OTP verification
func VerifyOTP(c *gin.Context) {
	var req models.OTPVerify
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	otpData, exists := utils.GetOTP(req.PhoneNumber)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No OTP request found for this phone number"})
		return
	}

	// Check if OTP has expired (5 minutes)
	if time.Since(otpData.CreatedAt) > 5*time.Minute {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
		return
	}

	// Verify OTP
	if otpData.OTP != req.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	// Check if user exists
	user, exists := utils.GetUserByPhone(req.PhoneNumber)
	var token string
	var err error
	
	if exists {
		// Generate JWT token for existing user
		token, err = utils.GenerateJWT(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
	} else {
		// For phone numbers without registered users, generate a temporary token
		// This token will be used for registration
		tempUser := models.User{
			ID:    "temp_" + req.PhoneNumber,
			Phone: req.PhoneNumber,
		}
		token, err = utils.GenerateJWT(tempUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
	}

	// Remove OTP from store after successful verification
	utils.DeleteOTP(req.PhoneNumber)
	
	response := gin.H{
		"message": "Phone number verified successfully",
		"token":   token,
	}
	
	c.JSON(http.StatusOK, response)
}
