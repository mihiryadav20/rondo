package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// TwilioClient handles Twilio API operations
type TwilioClient struct {
	Client     *twilio.RestClient
	AccountSID string
	AuthToken  string
	FromNumber string
}

// InitTwilio initializes and returns a new Twilio client
func InitTwilio() *TwilioClient {
	// Get credentials from environment variables
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")

	// Check if credentials are available
	if accountSid == "" || authToken == "" || fromNumber == "" {
		log.Println("Warning: One or more Twilio credentials are missing from environment variables")
		log.Println("Make sure TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, and TWILIO_FROM_NUMBER are set in your .env file")
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

// SendOTP sends an OTP via Twilio SMS
func (tc *TwilioClient) SendOTP(phoneNumber, otp string) error {
	// Create the message params
	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(tc.FromNumber)
	params.SetBody(fmt.Sprintf("Hello user, the verification code is: %s", otp))

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
