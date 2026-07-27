package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type AuthContext struct {
	UserID uuid.UUID
	Role   UserRole
}

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	DisplayName  string    `json:"displayName" db:"display_name"`
	Email        string    `json:"-" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         UserRole  `json:"role" db:"role"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	UpdatedAt    time.Time `json:"-" db:"updated_at"`
}

func (User) Table() string {
	return "users"
}
