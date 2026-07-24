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

func scanRecipe(row pgx.Row) (*models.Recipe, error) {
	var recipe models.Recipe
	err := row.Scan(
		&recipe.ID,
		&recipe.AuthorID,
		&recipe.Title,
		&recipe.Description,
		&recipe.Ingredients,
		&recipe.Instructions,
		&recipe.PrepTimeMin,
		&recipe.Servings,
		&recipe.Difficulty,
		&recipe.RatingAvg,
		&recipe.RatingCount,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *RecipeRepository) Create(ctx context.Context, recipe *models.Recipe) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO recipes (
			id,
			author_id,
			title,
			description,
			ingredients,
			instructions,
			prep_time_min,
			servings,
			difficulty
		 ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		recipe.ID,
		recipe.AuthorID,
		recipe.Title,
		recipe.Description,
		recipe.Ingredients,
		recipe.Instructions,
		recipe.PrepTimeMin,
		recipe.Servings,
		recipe.Difficulty,
	)
	return err
}

func (r *RecipeRepository) Get(ctx context.Context, id uuid.UUID) (*models.Recipe, error) {
	recipe, err := scanRecipe(r.db.QueryRow(ctx,
		`SELECT
       		id,
			author_id,
			title,
			description,
			ingredients,
			instructions,
			prep_time_min,
			servings,
			difficulty,
			rating_avg,
			rating_count,
			created_at,
			updated_at
     	 FROM recipes
     	 WHERE id = $1`,
		id,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecipeNotFound
		}
		return nil, err
	}
	return recipe, nil
}

func (r *RecipeRepository) List(ctx context.Context, limit, offset int) ([]*models.Recipe, error) {
	rows, err := r.db.Query(ctx,
		`SELECT
			id,
			author_id,
			title,
			description,
			ingredients,
			instructions,
			prep_time_min,
			servings,
			difficulty,
			rating_avg,
			rating_count,
			created_at,
			updated_at
		 FROM recipes
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2;`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var recipes []*models.Recipe
	for rows.Next() {
		recipe, err := scanRecipe(rows)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return recipes, nil
}

func (r *RecipeRepository) Update(ctx context.Context, recipe *models.Recipe) error {
	result, err := r.db.Exec(ctx,
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
		 WHERE id = $1`,
		recipe.ID,
		recipe.Title,
		recipe.Description,
		recipe.Ingredients,
		recipe.Instructions,
		recipe.PrepTimeMin,
		recipe.Servings,
		recipe.Difficulty,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRecipeNotFound
	}

	return nil
}

func (r *RecipeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM recipes WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrRecipeNotFound
	}
	return nil
}
