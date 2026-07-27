package services

import (
	"context"
	"errors"
	"strings"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/ports"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/repositories"
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

func (s *UserService) Register(ctx context.Context,
	displayName string,
	email string,
	password string,
) (*models.User, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	_, err := s.repository.GetByEmail(ctx, email)

	switch {
	case err == nil:
		return nil, ErrEmailAlreadyExists
	case errors.Is(err, repositories.ErrUserNotFound):
		// continue
	default:
		return nil, err
	}

	hash, err := s.hasher.Hash(password)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		DisplayName:  displayName,
		Email:        email,
		PasswordHash: hash,
		Role:         models.RoleUser,
	}

	err = s.repository.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.repository.GetByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	err = s.hasher.Compare(
		password,
		user.PasswordHash,
	)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(user.ID, user.Role)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.User, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *UserService) Update(
	ctx context.Context,
	userID uuid.UUID,
	user *models.User,
) error {
	//if user.ID != authUser.ID && authUser.Role != models.RoleAdmin
	if user.ID != userID {
		return ErrUnauthorized
	}

	user.Email = strings.ToLower(
		strings.TrimSpace(user.Email),
	)

	existing, err := s.repository.GetByEmail(ctx, user.Email)

	if err == nil {
		if existing.ID != user.ID {
			return ErrEmailAlreadyExists
		}
	} else if !errors.Is(err, repositories.ErrUserNotFound) {
		return err
	}

	return s.repository.Update(ctx, user)
}

func (s *UserService) Delete(
	ctx context.Context,
	userID uuid.UUID,
	id uuid.UUID,
) error {
	if id != userID {
		return ErrUnauthorized
	}

	return s.repository.Delete(ctx, id)
}
