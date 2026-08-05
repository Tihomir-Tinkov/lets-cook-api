package services

import (
	"context"
	"errors"
	"strings"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/dto"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/models"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/ports"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/repositories"
	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
)

type UserService struct {
	repository ports.UserRepository
	hasher     ports.PasswordHasher
	tokens     ports.TokenProvider
}

func NewUserService(repository ports.UserRepository, hasher ports.PasswordHasher, tokens ports.TokenProvider) *UserService {
	return &UserService{
		repository: repository,
		hasher:     hasher,
		tokens:     tokens,
	}
}

func mapUserToResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        user.Role,
	}
}

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	_, err := s.repository.GetByEmail(ctx, req.Email)

	switch {
	case err == nil:
		return nil, ErrEmailAlreadyExists
	case errors.Is(err, repositories.ErrUserNotFound):
		// continue
	default:
		return nil, err
	}

	hash, err := s.hasher.Hash(req.Password)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		DisplayName:  req.DisplayName,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         models.RoleUser,
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return mapUserToResponse(user), nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.repository.GetByEmail(ctx, req.Email)

	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	err = s.hasher.Compare(
		req.Password,
		user.PasswordHash,
	)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(user.ID, user.Role)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		User:  *mapUserToResponse(user),
		Token: token,
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.repository.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return mapUserToResponse(user), nil
}

func (s *UserService) Update(ctx context.Context, userID uuid.UUID, req dto.UserUpdateRequest) (*dto.UserResponse, error) {
	user, err := s.repository.GetByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	user.DisplayName = req.DisplayName
	user.Email = strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := s.repository.GetByEmail(ctx, user.Email)

	if err == nil {
		if existing.ID != user.ID {
			return nil, ErrEmailAlreadyExists
		}
	} else if !errors.Is(err, repositories.ErrUserNotFound) {
		return nil, err
	}

	if req.Password != "" {
		hash, err := s.hasher.Hash(req.Password)

		if err != nil {
			return nil, err
		}

		user.PasswordHash = hash
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return nil, err
	}

	return mapUserToResponse(user), nil
}

func (s *UserService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	if id != userID {
		return ErrUnauthorized
	}

	return s.repository.Delete(ctx, id)
}
