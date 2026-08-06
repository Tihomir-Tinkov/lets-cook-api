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
	ErrInvalidComment   = errors.New("invalid comment")
	ErrInvalidContent   = errors.New("comment text is required")
	ErrInvalidRating    = errors.New("rating must be between 1 and 5")
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

	if strings.TrimSpace(comment.Content) == "" {
		return ErrInvalidContent
	}

	if comment.Rating < 1 || comment.Rating > 5 {
		return ErrInvalidRating
	}

	return nil
}

func mapCommentToResponse(comment *models.Comment) *dto.CommentResponse {
	return &dto.CommentResponse{
		ID:       comment.ID,
		RecipeID: comment.RecipeID,
		AuthorID: comment.AuthorID,
		Content:  comment.Content,
		Rating:   comment.Rating,
	}
}

func (s *CommentService) Create(ctx context.Context, recipeID uuid.UUID, userID uuid.UUID, req dto.CommentCreateRequest) (*dto.CommentResponse, error) {
	comment := &models.Comment{
		RecipeID: recipeID,
		AuthorID: userID,
		Content:  req.Content,
		Rating:   req.Rating,
	}

	if err := validateComment(comment); err != nil {
		return nil, err
	}

	if err := s.repository.Create(ctx, comment); err != nil {
		return nil, err
	}

	return mapCommentToResponse(comment), nil
}

func (s *CommentService) GetByID(ctx context.Context, id uuid.UUID) (*dto.CommentResponse, error) {
	comment, err := s.repository.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return mapCommentToResponse(comment), nil
}

func (s *CommentService) GetByRecipeID(ctx context.Context, recipeID uuid.UUID) ([]dto.CommentResponse, error) {
	comments, err := s.repository.GetByRecipeID(ctx, recipeID)

	if err != nil {
		return nil, err
	}

	result := make([]dto.CommentResponse, 0, len(comments))

	for _, comment := range comments {
		result = append(result, *mapCommentToResponse(&comment))
	}

	return result, nil
}

func (s *CommentService) Update(ctx context.Context, recipeID uuid.UUID, userID uuid.UUID, commentID uuid.UUID, req dto.CommentUpdateRequest) (*dto.CommentResponse, error) {
	existing, err := s.repository.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if existing.RecipeID != recipeID {
		return nil, repositories.ErrCommentNotFound
	}

	if existing.AuthorID != userID {
		return nil, ErrCommentForbidden
	}

	existing.Content = req.Content
	existing.Rating = req.Rating

	if err := validateComment(existing); err != nil {
		return nil, err
	}

	if err := s.repository.Update(ctx, existing); err != nil {
		return nil, err
	}

	return mapCommentToResponse(existing), nil
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
