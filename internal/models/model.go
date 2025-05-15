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

type UserRole struct {
	UserID     uuid.UUID `json:"user_id"`
	RoleID     uuid.UUID `json:"role_id"`
	AssignedBy uuid.UUID `json:"assigned_by"`
}

type PermissionShort struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type PermissionDetails struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserShortInfo struct {
	UserName      string `json:"username"`
	Email         string `json:"email"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	EmailVerified bool   `json:"email_verified"`
	Active        bool   `json:"active"`
}
