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

func (r *RecipeRepository) Create(ctx context.Context, recipe *models.Recipe, images []models.RecipeImage) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
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

	if err != nil {
		return err
	}

	for i := range images {
		images[i].RecipeID = recipe.ID
	}

	for _, img := range images {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO recipe_images (
				recipe_id,
				image_id,
				display_order
			)
			VALUES ($1,$2,$3)`,
			img.RecipeID,
			img.ImageID,
			img.DisplayOrder,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *RecipeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Recipe, error) {
	row := r.db.QueryRow(
		ctx,
		`
		SELECT
			r.id,
			r.author_id,
			r.title,
			r.description,
			r.ingredients,
			r.instructions,
			r.prep_time_min,
			r.servings,
			r.difficulty,
			r.rating_avg,
			r.rating_count,
			r.created_at,
			r.updated_at,

			COALESCE(
				json_agg(
					json_build_object(
						'imageId', ri.image_id,
						'recipeId', ri.recipe_id,
						'displayOrder', ri.display_order
					)
					ORDER BY ri.display_order
				) FILTER (WHERE ri.image_id IS NOT NULL),
				'[]'
			) AS images

		FROM recipes r
		LEFT JOIN recipe_images ri
			ON ri.recipe_id = r.id
		WHERE r.id = $1
		GROUP BY r.id;
		`,
		id,
	)

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
		&recipe.Images,
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

func (r *RecipeRepository) Update(ctx context.Context, recipe *models.Recipe, images []models.RecipeImage) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
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

	// Delete images not present anymore
	imageIDs := make([]uuid.UUID, len(images))

	for i, img := range images {
		imageIDs[i] = img.ImageID
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM recipe_images
      	WHERE recipe_id=$1`,
		recipe.ID,
	)
	if err != nil {
		return err
	}

	for _, img := range images {
		_, err = tx.Exec(ctx,
			`INSERT INTO recipe_images
				(recipe_id,image_id,display_order)
			VALUES ($1,$2,$3)`,
			recipe.ID,
			img.ImageID,
			img.DisplayOrder,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
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
