package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Verified    bool   `json:"verified"`
}

type OTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type OTPVerify struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
}

// OTPData stores OTP information
type OTPData struct {
	OTP       string
	CreatedAt time.Time
}

// TwilioClient handles Twilio API operations
type TwilioClient struct {
	Client     *twilio.RestClient
	AccountSID string
	AuthToken  string
	FromNumber string
}

// Global variables
var (
	otpStore     = make(map[string]OTPData)
	otpStoreLock sync.RWMutex
	twilioClient *TwilioClient
)

var users = []User{
	{ID: "1", Name: "John Doe", Email: "john@example.com", PhoneNumber: "+911234567890", Verified: false},
	{ID: "2", Name: "Jane Smith", Email: "jane@example.com", PhoneNumber: "+911234567891", Verified: false},
}

// Helper functions for OTP
func generateOTP() string {
	buffer := make([]byte, 3)
	rand.Read(buffer)
	num := (int(buffer[0])*256*256 + int(buffer[1])*256 + int(buffer[2])) % 1000000
	return fmt.Sprintf("%06d", num)
}

// Initialize Twilio client
func initTwilio() *TwilioClient {
	// In production, use environment variables
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID") // Replace with your Twilio Account SID
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")   // Replace with your Twilio Auth Token
	fromNumber := os.Getenv("TWILIO_FROM_NUMBER") // Replace with your Twilio phone number

	// For development, you can hardcode these values (not recommended for production)
	if accountSid == "" {
		accountSid = "YOUR_TWILIO_ACCOUNT_SID" // Replace with your Twilio Account SID
	}
	if authToken == "" {
		authToken = "YOUR_TWILIO_AUTH_TOKEN" // Replace with your Twilio Auth Token
	}
	if fromNumber == "" {
		fromNumber = "YOUR_TWILIO_PHONE_NUMBER" // Replace with your Twilio phone number
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return &TwilioClient{
		Client:     client,
		AccountSID: accountSid,
		AuthToken:  authToken,
		FromNumber: fromNumber,
	}
}

// Send OTP via Twilio SMS
func (tc *TwilioClient) sendOTP(phoneNumber, otp string) error {
	// Create the message params
	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(tc.FromNumber)
	params.SetBody(fmt.Sprintf("Your verification code is: %s", otp))

	// Send the message
	_, err := tc.Client.Api.CreateMessage(params)
	if err != nil {
		fmt.Printf("Error sending SMS: %s\n", err.Error())
		return err
	}

	// For development, also print to console
	fmt.Printf("Sending OTP %s to %s via Twilio\n", otp, phoneNumber)
	return nil
}

// Request OTP handler
func requestOTP(c *gin.Context) {
	var req OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	otp := generateOTP()
	otpStoreLock.Lock()
	otpStore[req.PhoneNumber] = OTPData{
		OTP:       otp,
		CreatedAt: time.Now(),
	}
	otpStoreLock.Unlock()

	// Send OTP via Twilio
	if err := twilioClient.sendOTP(req.PhoneNumber, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// Verify OTP handler
func verifyOTP(c *gin.Context) {
	var req OTPVerify
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	otpStoreLock.RLock()
	otpData, exists := otpStore[req.PhoneNumber]
	otpStoreLock.RUnlock()

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

	// Update user verification status if phone number matches
	userFound := false
	for i, user := range users {
		if user.PhoneNumber == req.PhoneNumber {
			users[i].Verified = true
			userFound = true
			break
		}
	}

	if !userFound {
		// In a real app, you might want to create a new user or handle this differently
		c.JSON(http.StatusNotFound, gin.H{"error": "No user found with this phone number"})
		return
	}

	// Remove OTP from store after successful verification
	otpStoreLock.Lock()
	delete(otpStore, req.PhoneNumber)
	otpStoreLock.Unlock()
	c.JSON(http.StatusOK, gin.H{"message": "Phone number verified successfully"})
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Initialize Twilio client
	twilioClient = initTwilio()

	r := gin.Default()

	// Basic health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Get all users
	r.GET("/users", getUsers)

	// Get user by ID
	r.GET("/users/:id", getUserByID)

	// Create a new user
	r.POST("/users", createUser)

	// OTP endpoints
	r.POST("/auth/otp/request", requestOTP)
	r.POST("/auth/otp/verify", verifyOTP)

	// Start the server
	r.Run(":8080")
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real app, you would generate a unique ID and save to a database
	users = append(users, newUser)
	c.JSON(http.StatusCreated, newUser)
}
