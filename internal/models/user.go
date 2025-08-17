package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewUser creates a new user with validation
func NewUser(name, email, phoneNumber string) (*User, error) {
	if name == "" || email == "" || phoneNumber == "" {
		return nil, ErrInvalidUserData
	}

	return &User{
		ID:          uuid.New().String(),
		Name:        name,
		Email:       email,
		PhoneNumber: phoneNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(name, email, phoneNumber string) error {
	if name == "" || email == "" || phoneNumber == "" {
		return ErrInvalidUserData
	}

	u.Name = name
	u.Email = email
	u.PhoneNumber = phoneNumber
	u.UpdatedAt = time.Now()
	return nil
}
