package models

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID           uuid.UUID `json:"id" db:"id"`
	AuthorID     uuid.UUID `json:"authorId" db:"author_id"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	Ingredients  string    `json:"ingredients" db:"ingredients"`
	Instructions string    `json:"instructions" db:"instructions"`
	PrepTimeMin  int       `json:"prepTimeMin" db:"prep_time_min"`
	Servings     int       `json:"servings" db:"servings"`
	Difficulty   int       `json:"difficulty" db:"difficulty"`
	RatingAvg    float64   `json:"ratingAvg" db:"rating_avg"`
	RatingCount  int       `json:"ratingCount" db:"rating_count"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	UpdatedAt    time.Time `json:"-" db:"updated_at"`
}

func (Recipe) Table() string {
	return "recipes"
}
