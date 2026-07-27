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
	ErrInvalidIngredients  = errors.New("ingredients are required")
	ErrInvalidInstructions = errors.New("instructions are required")
	ErrInvalidPrepTime     = errors.New("prep time must be greater than zero")
	ErrInvalidServings     = errors.New("servings must be greater than zero")
	ErrInvalidDifficulty   = errors.New("difficulty must be between 1 and 10")
	ErrForbidden           = errors.New("recipe access forbidden")
)

type RecipeService struct {
	recipeRepository ports.RecipeRepository
}

func NewRecipeService(recipeRepository ports.RecipeRepository) *RecipeService {
	return &RecipeService{
		recipeRepository: recipeRepository,
	}
}

func validateRecipe(recipe *models.Recipe) error {
	if recipe == nil {
		return ErrInvalidRecipe
	}

	if strings.TrimSpace(recipe.Title) == "" {
		return ErrInvalidTitle
	}

	if strings.TrimSpace(recipe.Ingredients) == "" {
		return ErrInvalidIngredients
	}

	if strings.TrimSpace(recipe.Instructions) == "" {
		return ErrInvalidInstructions
	}

	if recipe.PrepTimeMin <= 0 {
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

	return s.recipeRepository.Create(ctx, recipe)
}

func (s *RecipeService) GetByID(ctx context.Context, id uuid.UUID) (*models.Recipe, error) {
	return s.recipeRepository.GetByID(ctx, id)
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

	return s.recipeRepository.List(ctx, limit, offset)
}

func (s *RecipeService) Update(ctx context.Context, userID uuid.UUID, recipe *models.Recipe) error {
	if recipe.ID == uuid.Nil {
		return ErrInvalidRecipe
	}

	if err := validateRecipe(recipe); err != nil {
		return err
	}

	existing, err := s.recipeRepository.GetByID(ctx, recipe.ID)
	if err != nil {
		return err
	}

	if existing.AuthorID != userID {
		return ErrForbidden
	}

	return s.recipeRepository.Update(ctx, recipe)
}

func (s *RecipeService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	recipe, err := s.recipeRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if recipe.AuthorID != userID {
		return ErrForbidden
	}

	return s.recipeRepository.Delete(ctx, id)
}
