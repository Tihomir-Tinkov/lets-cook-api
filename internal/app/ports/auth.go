package ports

import (
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/models"
	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password string, hash string) error
}

type TokenProvider interface {
	Generate(userID uuid.UUID, role models.UserRole) (string, error)
	Validate(token string) (*models.AuthContext, error)
}
