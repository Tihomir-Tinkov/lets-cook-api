package dto

import (
	"github.com/google/uuid"
)

type CommentCreateRequest struct {
	Content string `json:"content"`
	Rating  int    `json:"rating"`
}

type CommentUpdateRequest struct {
	Content string `json:"content"`
	Rating  int    `json:"rating"`
}

type CommentResponse struct {
	ID       uuid.UUID `json:"id" db:"id"`
	RecipeID uuid.UUID `json:"recipeId" db:"recipe_id"`
	AuthorID uuid.UUID `json:"authorId" db:"author_id"`
	Content  string    `json:"content" db:"content"`
	Rating   int       `json:"rating" db:"rating"`
}
