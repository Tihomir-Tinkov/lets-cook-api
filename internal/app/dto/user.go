package dto

import (
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/models"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type UserResponse struct {
	ID          uuid.UUID       `json:"id"`
	DisplayName string          `json:"displayName"`
	Email       string          `json:"email"`
	Role        models.UserRole `json:"role" db:"role"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
