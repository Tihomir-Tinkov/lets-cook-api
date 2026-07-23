package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
)

type UserService interface {
	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (models.User, error)
}
