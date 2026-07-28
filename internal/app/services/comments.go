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
	ErrInvalidComment   = errors.New("invalid comment")
	ErrInvalidBody      = errors.New("comment body is required")
	ErrInvalidRating    = errors.New("rating must be between 1 and 10")
	ErrCommentForbidden = errors.New("comment access forbidden")
)

type CommentService struct {
	repository ports.CommentRepository
}

func NewCommentService(repository ports.CommentRepository) *CommentService {
	return &CommentService{
		repository: repository,
	}
}

func validateComment(comment *models.Comment) error {
	if comment == nil {
		return ErrInvalidComment
	}

	if strings.TrimSpace(comment.Body) == "" {
		return ErrInvalidBody
	}

	if comment.Rating < 1 || comment.Rating > 10 {
		return ErrInvalidRating
	}

	return nil
}

func (s *CommentService) Create(ctx context.Context, comment *models.Comment) error {
	if err := validateComment(comment); err != nil {
		return err
	}

	return s.repository.Create(ctx, comment)
}

func (s *CommentService) GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *CommentService) GetByRecipeID(ctx context.Context, recipeID uuid.UUID) ([]models.Comment, error) {
	return s.repository.GetByRecipeID(ctx, recipeID)
}

func (s *CommentService) Update(ctx context.Context, recipeID uuid.UUID, userID uuid.UUID, comment *models.Comment) error {
	if err := validateComment(comment); err != nil {
		return err
	}

	if comment.ID == uuid.Nil {
		return ErrInvalidComment
	}

	existing, err := s.repository.GetByID(ctx, comment.ID)
	if err != nil {
		return err
	}

	if existing.RecipeID != recipeID {
		return repositories.ErrCommentNotFound
	}

	if existing.AuthorID != userID {
		return ErrCommentForbidden
	}

	return s.repository.Update(ctx, comment)
}

func (s *CommentService) Delete(ctx context.Context, recipeID uuid.UUID, userID uuid.UUID, id uuid.UUID) error {
	comment, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if comment.RecipeID != recipeID {
		return repositories.ErrCommentNotFound
	}

	if comment.AuthorID != userID {
		return ErrCommentForbidden
	}

	return s.repository.Delete(ctx, id)
}
