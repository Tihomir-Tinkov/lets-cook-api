package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/dto"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/models"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/ports"
	"github.com/google/uuid"
)

var (
	ErrInvalidRecipe       = errors.New("invalid recipe")
	ErrInvalidTitle        = errors.New("title is required")
	ErrInvalidDescription  = errors.New("description is required")
	ErrInvalidIngredients  = errors.New("ingredients are required")
	ErrInvalidInstructions = errors.New("instructions are required")
	ErrInvalidPrepTime     = errors.New("prep time cannot be negative")
	ErrInvalidDifficulty   = errors.New("difficulty must be easy, medium, or hard")
	ErrInvalidServings     = errors.New("servings must be greater than zero")
	ErrInvalidImages       = errors.New("recipe must have at least one image")
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

	switch recipe.Difficulty {
	case models.Easy, models.Medium, models.Hard:
		// valid
	default:
		return ErrInvalidDifficulty
	}

	if recipe.Servings <= 0 {
		return ErrInvalidServings
	}

	if len(recipe.Images) == 0 {
		return ErrInvalidRecipe
	}

	return nil
}

func mapRecipeImages(images []models.RecipeImage) []dto.RecipeImageResponse {
	result := make([]dto.RecipeImageResponse, 0, len(images))

	for _, img := range images {
		result = append(result, dto.RecipeImageResponse{
			ImageID:      img.ImageID,
			URL:          fmt.Sprintf("/images/%s", img.ImageID),
			DisplayOrder: img.DisplayOrder,
		})
	}

	return result
}

func mapRecipeToResponse(recipe *models.Recipe) *dto.RecipeResponse {
	return &dto.RecipeResponse{
		ID:           recipe.ID,
		AuthorID:     recipe.AuthorID,
		Title:        recipe.Title,
		Description:  recipe.Description,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		PrepTimeMin:  recipe.PrepTimeMin,
		Difficulty:   recipe.Difficulty,
		Servings:     recipe.Servings,
		RatingAvg:    recipe.RatingAvg,
		RatingCount:  recipe.RatingCount,
		Images:       mapRecipeImages(recipe.Images),
	}
}

func mapRecipeToSummary(recipe models.Recipe) dto.RecipeSummary {
	return dto.RecipeSummary{
		ID:          recipe.ID,
		AuthorID:    recipe.AuthorID,
		Title:       recipe.Title,
		PrepTimeMin: recipe.PrepTimeMin,
		Difficulty:  recipe.Difficulty,
		Servings:    recipe.Servings,
		RatingAvg:   recipe.RatingAvg,
		RatingCount: recipe.RatingCount,
		Thumbnail:   mapRecipeImages(recipe.Images),
	}
}

func (s *RecipeService) Create(ctx context.Context, userID uuid.UUID, req dto.RecipeCreateRequest) (*dto.RecipeResponse, error) {
	recipe := &models.Recipe{
		AuthorID:     userID,
		Title:        req.Title,
		Description:  req.Description,
		Ingredients:  req.Ingredients,
		Instructions: req.Instructions,
		PrepTimeMin:  req.PrepTimeMin,
		Difficulty:   req.Difficulty,
		Servings:     req.Servings,
	}

	for _, img := range req.Images {
		recipe.Images = append(recipe.Images, models.RecipeImage{
			ImageID:      img.ImageID,
			DisplayOrder: img.DisplayOrder,
		})
	}

	if err := validateRecipe(recipe); err != nil {
		return nil, err
	}

	err := s.repository.Create(ctx, recipe)

	if err != nil {
		return nil, err
	}

	return mapRecipeToResponse(recipe), nil
}

func (s *RecipeService) GetByID(ctx context.Context, id uuid.UUID) (*dto.RecipeResponse, error) {
	recipe, err := s.repository.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return mapRecipeToResponse(recipe), nil
}

func (s *RecipeService) List(ctx context.Context, limit, offset int) (*dto.RecipeListResponse, error) {
	if limit <= 0 {
		limit = 12
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	list, err := s.repository.List(ctx, limit, offset)

	if err != nil {
		return nil, err
	}

	recipes := make([]dto.RecipeSummary, 0, len(list))

	for _, recipe := range list {
		recipes = append(recipes, mapRecipeToSummary(recipe))
	}

	totalItems, err := s.repository.Count(ctx)

	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	return &dto.RecipeListResponse{
		Limit:      limit,
		Page:       offset/limit + 1,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Recipes:    recipes,
	}, nil
}

func (s *RecipeService) Update(ctx context.Context, recipeID uuid.UUID, userID uuid.UUID, req dto.RecipeUpdateRequest) (*dto.RecipeResponse, error) {
	existing, err := s.repository.GetByID(ctx, recipeID)
	if err != nil {
		return nil, err
	}

	if existing.AuthorID != userID {
		return nil, ErrRecipeForbidden
	}

	existing.Title = req.Title
	existing.Description = req.Description
	existing.Ingredients = req.Ingredients
	existing.Instructions = req.Instructions
	existing.PrepTimeMin = req.PrepTimeMin
	existing.Difficulty = req.Difficulty
	existing.Servings = req.Servings

	existing.Images = make([]models.RecipeImage, len(req.Images))
	for i, img := range req.Images {
		existing.Images[i] = models.RecipeImage{
			ImageID:      img.ImageID,
			DisplayOrder: img.DisplayOrder,
		}
	}

	if err := validateRecipe(existing); err != nil {
		return nil, err
	}

	if err := s.repository.Update(ctx, existing); err != nil {
		return nil, err
	}

	return mapRecipeToResponse(existing), nil
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
