package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	Username          string     `json:"username" db:"username"`
	Email             string     `json:"email" db:"email"`
	PasswordHash      string     `json:"password_hash" db:"password_hash"`
	FirstName         *string    `json:"first_name,omitempty" db:"first_name"`
	LastName          *string    `json:"last_name,omitempty" db:"last_name"`
	EmailVerified     bool       `json:"email_verified" db:"email_verified"`
	UserType          string     `json:"user_type" db:"user_type"`
	VerificationToken *string    `json:"verification_token,omitempty" db:"verification_token"`
	TokenExpiry       *time.Time `json:"token_expiry,omitempty" db:"token_expiry"`
	DeletionRequested bool       `json:"deletion_requested" db:"deletion_requested"`
	Active            bool       `json:"active" db:"active"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}
