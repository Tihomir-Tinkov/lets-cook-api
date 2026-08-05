package dto

import (
	"github.com/google/uuid"
)

type CommentCreateRequest struct {
	Body   string `json:"body"`
	Rating int    `json:"rating"`
}

type CommentUpdateRequest struct {
	Body   string `json:"body"`
	Rating int    `json:"rating"`
}

type CommentResponse struct {
	ID       uuid.UUID `json:"id" db:"id"`
	RecipeID uuid.UUID `json:"recipeId" db:"recipe_id"`
	AuthorID uuid.UUID `json:"authorId" db:"author_id"`
	Body     string    `json:"body" db:"body"`
	Rating   int       `json:"rating" db:"rating"`
}
