package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	RecipeID  uuid.UUID `json:"recipeId" db:"recipe_id"`
	AuthorID  uuid.UUID `json:"authorId" db:"author_id"`
	Content   string    `json:"content" db:"content"`
	Rating    int       `json:"rating" db:"rating"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

func (Comment) Table() string {
	return "comments"
}
