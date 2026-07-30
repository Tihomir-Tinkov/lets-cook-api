package services

import (
	"context"
	"errors"
	"strings"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/ports"
	"github.com/google/uuid"
)

var (
	ErrInvalidRecipe       = errors.New("invalid recipe")
	ErrInvalidTitle        = errors.New("title is required")
	ErrInvalidDescription  = errors.New("description is required")
	ErrInvalidIngredients  = errors.New("ingredients are required")
	ErrInvalidInstructions = errors.New("instructions are required")
	ErrInvalidPrepTime     = errors.New("prep time cannot be negative")
	ErrInvalidServings     = errors.New("servings must be greater than zero")
	ErrInvalidDifficulty   = errors.New("difficulty must be between 1 and 10")
	ErrRecipeForbidden     = errors.New("recipe access forbidden")
)

type RecipeService struct {
	repository ports.RecipeRepository
}

func NewRecipeService(repository ports.RecipeRepository) *RecipeService {
	return &RecipeService{
		repository: repository,
	}
}

func validateRecipe(recipe *models.Recipe) error {
	if recipe == nil {
		return ErrInvalidRecipe
	}

	if strings.TrimSpace(recipe.Title) == "" {
		return ErrInvalidTitle
	}

	if strings.TrimSpace(recipe.Description) == "" {
		return ErrInvalidDescription
	}

	if strings.TrimSpace(recipe.Ingredients) == "" {
		return ErrInvalidIngredients
	}

	if strings.TrimSpace(recipe.Instructions) == "" {
		return ErrInvalidInstructions
	}

	if recipe.PrepTimeMin < 0 {
		return ErrInvalidPrepTime
	}

	if recipe.Servings <= 0 {
		return ErrInvalidServings
	}

	if recipe.Difficulty < 1 || recipe.Difficulty > 10 {
		return ErrInvalidDifficulty
	}

	return nil
}

func (s *RecipeService) Create(ctx context.Context, recipe *models.Recipe) error {
	if err := validateRecipe(recipe); err != nil {
		return err
	}

	if len(recipe.Images) == 0 {
		return ErrInvalidRecipe
	}

	return s.repository.Create(ctx, recipe, recipe.Images)
}

func (s *RecipeService) GetByID(ctx context.Context, id uuid.UUID) (*models.Recipe, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *RecipeService) List(ctx context.Context, limit, offset int) ([]models.Recipe, error) {
	if limit <= 0 {
		limit = 12
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	return s.repository.List(ctx, limit, offset)
}

func (s *RecipeService) Update(ctx context.Context, userID uuid.UUID, recipe *models.Recipe) error {
	if err := validateRecipe(recipe); err != nil {
		return err
	}

	if recipe.ID == uuid.Nil {
		return ErrInvalidRecipe
	}

	if len(recipe.Images) == 0 {
		return ErrInvalidRecipe
	}

	existing, err := s.repository.GetByID(ctx, recipe.ID)
	if err != nil {
		return err
	}

	if existing.AuthorID != userID {
		return ErrRecipeForbidden
	}

	return s.repository.Update(ctx, recipe, recipe.Images)
}

func (s *RecipeService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	recipe, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if recipe.AuthorID != userID {
		return ErrRecipeForbidden
	}

	return s.repository.Delete(ctx, id)
}
