package ports

import (
	"context"
	"io"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RecipeRepository interface {
	Create(ctx context.Context, recipe *models.Recipe) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Recipe, error)
	List(ctx context.Context, limit, offset int) ([]models.Recipe, error)
	Update(ctx context.Context, recipe *models.Recipe) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CommentRepository interface {
	Create(ctx context.Context, recipe *models.Comment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error)
	GetByRecipeID(ctx context.Context, recipeID uuid.UUID) ([]models.Comment, error)
	Update(ctx context.Context, recipe *models.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ImageRepository interface {
	Create(ctx context.Context, image *models.Image) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Image, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type FileStorage interface {
	Save(ctx context.Context, id uuid.UUID, extension string, r io.Reader) error
	Open(ctx context.Context, id uuid.UUID, extension string) (io.ReadCloser, error)
	Delete(ctx context.Context, id uuid.UUID, extension string) error
}
