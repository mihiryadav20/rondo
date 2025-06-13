package models

import "time"

// OTPRequest represents the request to send an OTP
type OTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// OTPVerify represents the request to verify an OTP
type OTPVerify struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
}

// OTPData stores OTP information
type OTPData struct {
	OTP       string
	CreatedAt time.Time
}
