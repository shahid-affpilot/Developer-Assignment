package models

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"max=50"`
	LastName  string `json:"last_name" validate:"max=50"`
}

type RegisterResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Message  string    `json:"message"`
}

type User struct {
	ID                uuid.UUID  `json:"id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	EmailVerified     bool       `json:"email_verified"`
	UserType          string     `json:"user_type"`
	VerificationToken string     `json:"-"`
	TokenExpiry       *time.Time `json:"-"`
	DeletionRequested bool       `json:"deletion_requested"`
	Active            bool       `json:"active"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
