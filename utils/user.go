package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	
	"rondo/models"
)

// UserStore is a simple in-memory storage for users
var (
	UserStore     = make(map[string]models.User) // Phone number -> User
	UserStoreLock sync.RWMutex
)

// CreateUser creates a new user in the store
func CreateUser(req models.UserRegistrationRequest, phoneNumber string) (models.User, error) {
	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		return models.User{}, fmt.Errorf("invalid date format: %v", err)
	}
	
	// Generate a new UUID for the user
	id := uuid.New().String()
	
	// Create user
	now := time.Now()
	user := models.User{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DOB:       dob,
		Phone:     phoneNumber,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	// Store user
	UserStoreLock.Lock()
	defer UserStoreLock.Unlock()
	
	// Check if phone number already exists
	for _, existingUser := range UserStore {
		if existingUser.Phone == phoneNumber {
			return models.User{}, fmt.Errorf("phone number already registered")
		}
	}
	
	UserStore[phoneNumber] = user
	return user, nil
}

// GetUserByPhone retrieves a user by phone number
func GetUserByPhone(phone string) (models.User, bool) {
	UserStoreLock.RLock()
	defer UserStoreLock.RUnlock()
	
	user, exists := UserStore[phone]
	return user, exists
}
