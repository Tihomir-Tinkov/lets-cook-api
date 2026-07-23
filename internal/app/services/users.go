package services

import (
	"context"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/ports"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/repositories"
	"github.com/google/uuid"
)

var _ ports.UserService = (*UserService)(nil)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(
	userRepository *repositories.UserRepository,
) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (models.User, error) {
	return s.userRepository.Find(ctx, id)
}
