package services

import (
	"context"
	"errors"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/repositories"
	"github.com/google/uuid"
)

// TODO: style consistency and if it works

var (
	ErrInvalidTitle        = errors.New("title is required")
	ErrInvalidIngredients  = errors.New("ingredients are required")
	ErrInvalidInstructions = errors.New("instructions are required")
	ErrInvalidPrepTime     = errors.New("prep time must be greater than zero")
	ErrInvalidServings     = errors.New("servings must be greater than zero")
	ErrInvalidDifficulty   = errors.New("difficulty must be between 1 and 10")
	ErrForbidden           = errors.New("forbidden")
)

type RecipeService struct {
	recipeRepository *repositories.RecipeRepository
}

func NewRecipeService(recipeRepository *repositories.RecipeRepository) *RecipeService {
	return &RecipeService{
		recipeRepository: recipeRepository,
	}
}

func validateRecipe(recipe *models.Recipe) error {
	if recipe.Title == "" {
		return ErrInvalidTitle
	}

	if recipe.Ingredients == "" {
		return ErrInvalidIngredients
	}

	if recipe.Instructions == "" {
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

	if recipe.ID == uuid.Nil {
		recipe.ID = uuid.New()
	}

	return s.recipeRepository.Create(ctx, recipe)
}

func (s *RecipeService) Get(ctx context.Context, id uuid.UUID) (*models.Recipe, error) {
	return s.recipeRepository.Get(ctx, id)
}

func (s *RecipeService) List(ctx context.Context, limit, offset int) ([]*models.Recipe, error) {
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
	existing, err := s.recipeRepository.Get(ctx, recipe.ID)
	if err != nil {
		return err
	}

	if existing.AuthorID != userID {
		return ErrForbidden
	}

	if err := validateRecipe(recipe); err != nil {
		return err
	}

	return s.recipeRepository.Update(ctx, recipe)
}

func (s *RecipeService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	recipe, err := s.recipeRepository.Get(ctx, id)
	if err != nil {
		return err
	}

	if recipe.AuthorID != userID {
		return ErrForbidden
	}

	return s.recipeRepository.Delete(ctx, id)
}
