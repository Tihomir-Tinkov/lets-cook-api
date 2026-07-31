package models

import (
	"time"

	"github.com/google/uuid"
)

type RecipeImage struct {
	ImageID      uuid.UUID `json:"imageId" db:"image_id"`
	RecipeID     uuid.UUID `json:"-" db:"recipe_id"`
	DisplayOrder int       `json:"displayOrder" db:"display_order"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
}

func (RecipeImage) Table() string {
	return "recipe_images"
}
