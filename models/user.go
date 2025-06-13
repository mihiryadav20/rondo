package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id,omitempty"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	DOB       time.Time `json:"dob" binding:"required"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// UserRegistrationRequest represents the request to register a new user
type UserRegistrationRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	DOB       string `json:"dob" binding:"required"` // Format: YYYY-MM-DD
	// Phone number comes from the JWT token
}

// UserResponse represents the response after user registration
type UserResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	DOB       time.Time `json:"dob"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}
