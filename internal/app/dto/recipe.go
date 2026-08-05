package dto

import (
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/models"
	"github.com/google/uuid"
)

type RecipeCreateRequest struct {
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Ingredients  string               `json:"ingredients"`
	Instructions string               `json:"instructions"`
	PrepTimeMin  int                  `json:"timeMinutes"`
	Difficulty   models.Difficulty    `json:"difficulty"`
	Servings     int                  `json:"servings"`
	Images       []RecipeImageRequest `json:"images"`
}

type RecipeUpdateRequest struct {
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Ingredients  string               `json:"ingredients"`
	Instructions string               `json:"instructions"`
	PrepTimeMin  int                  `json:"timeMinutes"`
	Difficulty   models.Difficulty    `json:"difficulty"`
	Servings     int                  `json:"servings"`
	Images       []RecipeImageRequest `json:"images"`
}

type RecipeResponse struct {
	ID           uuid.UUID             `json:"id" db:"id"`
	AuthorID     uuid.UUID             `json:"authorId" db:"author_id"`
	Title        string                `json:"title" db:"title"`
	Description  string                `json:"description" db:"description"`
	Ingredients  string                `json:"ingredients" db:"ingredients"`
	Instructions string                `json:"instructions" db:"instructions"`
	PrepTimeMin  int                   `json:"timeMinutes" db:"prep_time_min"`
	Difficulty   models.Difficulty     `json:"difficulty" db:"difficulty"`
	Servings     int                   `json:"servings" db:"servings"`
	RatingAvg    float64               `json:"ratingAvg" db:"rating_avg"`
	RatingCount  int                   `json:"ratingCount" db:"rating_count"`
	Images       []RecipeImageResponse `json:"images"`
}

type RecipeSummary struct {
	ID          uuid.UUID             `json:"id" db:"id"`
	AuthorID    uuid.UUID             `json:"authorId" db:"author_id"`
	Title       string                `json:"title" db:"title"`
	PrepTimeMin int                   `json:"timeMinutes" db:"prep_time_min"`
	Difficulty  models.Difficulty     `json:"difficulty" db:"difficulty"`
	Servings    int                   `json:"servings" db:"servings"`
	RatingAvg   float64               `json:"ratingAvg" db:"rating_avg"`
	RatingCount int                   `json:"ratingCount" db:"rating_count"`
	Thumbnail   []RecipeImageResponse `json:"images"`
}

type RecipeListResponse struct {
	Limit      int             `json:"limit"`
	Page       int             `json:"page"`
	TotalItems int             `json:"totalItems"`
	TotalPages int             `json:"totalPages"`
	Recipes    []RecipeSummary `json:"recipes"`
}
