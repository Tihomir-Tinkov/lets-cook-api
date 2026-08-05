package dto

import (
	"github.com/google/uuid"
)

type RecipeImageRequest struct {
	ImageID      uuid.UUID `json:"imageId"`
	DisplayOrder int       `json:"displayOrder"`
}

type RecipeImageResponse struct {
	ImageID      uuid.UUID `json:"imageId"`
	URL          string    `json:"url"`
	DisplayOrder int       `json:"displayOrder"`
}
