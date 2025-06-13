package utils

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"rondo/models"
)

// Global variables for OTP storage
var (
	OtpStore     = make(map[string]models.OTPData)
	OtpStoreLock sync.RWMutex
)

// GenerateOTP generates a random 6-digit OTP
func GenerateOTP() string {
	buffer := make([]byte, 3)
	rand.Read(buffer)
	num := (int(buffer[0])*256*256 + int(buffer[1])*256 + int(buffer[2])) % 1000000
	return fmt.Sprintf("%06d", num)
}

// StoreOTP stores an OTP for a phone number
func StoreOTP(phoneNumber, otp string) {
	OtpStoreLock.Lock()
	defer OtpStoreLock.Unlock()
	
	OtpStore[phoneNumber] = models.OTPData{
		OTP:       otp,
		CreatedAt: time.Now(),
	}
}

// GetOTP retrieves an OTP for a phone number
func GetOTP(phoneNumber string) (models.OTPData, bool) {
	OtpStoreLock.RLock()
	defer OtpStoreLock.RUnlock()
	
	otpData, exists := OtpStore[phoneNumber]
	return otpData, exists
}

// DeleteOTP removes an OTP from the store
func DeleteOTP(phoneNumber string) {
	OtpStoreLock.Lock()
	defer OtpStoreLock.Unlock()
	
	delete(OtpStore, phoneNumber)
}
