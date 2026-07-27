package repositories

import (
	"context"
	"errors"
	"reflect"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRecipeNotFound = errors.New("recipe not found")

type RecipeRepository struct {
	db *pgxpool.Pool
}

func NewRecipeRepository(db *pgxpool.Pool) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) DB() *pgxpool.Pool {
	return r.db
}

func (r *RecipeRepository) Model() *models.Recipe {
	return &models.Recipe{}
}

func (r *RecipeRepository) Type() reflect.Type {
	return reflect.TypeOf(models.Recipe{})
}

func (r *RecipeRepository) Create(ctx context.Context, recipe *models.Recipe) error {
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO recipes (
			author_id,
			title,
			description,
			ingredients,
			instructions,
			prep_time_min,
			servings,
			difficulty
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, rating_avg, rating_count, created_at, updated_at`,
		recipe.AuthorID,
		recipe.Title,
		recipe.Description,
		recipe.Ingredients,
		recipe.Instructions,
		recipe.PrepTimeMin,
		recipe.Servings,
		recipe.Difficulty,
	).Scan(
		&recipe.ID,
		&recipe.RatingAvg,
		&recipe.RatingCount,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
	)

	return err
}

func (r *RecipeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Recipe, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT *
		 FROM recipes
		 WHERE id = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipe, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByName[models.Recipe],
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecipeNotFound
		}

		return nil, err
	}

	return &recipe, nil
}

func (r *RecipeRepository) List(ctx context.Context, limit, offset int) ([]models.Recipe, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT *
		 FROM recipes
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipes, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[models.Recipe],
	)

	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (r *RecipeRepository) Update(ctx context.Context, recipe *models.Recipe) error {
	err := r.db.QueryRow(
		ctx,
		`UPDATE recipes
		 SET
			title = $2,
			description = $3,
			ingredients = $4,
			instructions = $5,
			prep_time_min = $6,
			servings = $7,
			difficulty = $8,
			updated_at = NOW()
		 WHERE id = $1
		 RETURNING updated_at`,
		recipe.ID,
		recipe.Title,
		recipe.Description,
		recipe.Ingredients,
		recipe.Instructions,
		recipe.PrepTimeMin,
		recipe.Servings,
		recipe.Difficulty,
	).Scan(&recipe.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRecipeNotFound
		}

		return err
	}

	return nil
}

func (r *RecipeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM recipes
		 WHERE id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRecipeNotFound
	}

	return nil
}
