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

var ErrCommentNotFound = errors.New("comment not found")

type CommentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) DB() *pgxpool.Pool {
	return r.db
}

func (r *CommentRepository) Model() *models.Comment {
	return &models.Comment{}
}

func (r *CommentRepository) Type() reflect.Type {
	return reflect.TypeOf(models.Comment{})
}

// updated_at = NOW() ?
func updateRecipeStats(ctx context.Context, recipeID uuid.UUID, tx pgx.Tx) error {
	_, err := tx.Exec(ctx,
		`UPDATE recipes r
		 SET
			rating_count = s.rating_count,
			rating_avg   = s.rating_avg
		 FROM (
			SELECT
				COUNT(*)::integer AS rating_count,
				COALESCE(AVG(rating), 0)::numeric(3,1) AS rating_avg
			FROM comments
			WHERE recipe_id = $1
		 ) s
		 WHERE r.id = $1`,
		recipeID,
	)

	return err
}

func (r *CommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			id,
    		recipe_id,
			author_id,
			body,
			rating,
			created_at,
			updated_at
		 FROM comments
		 WHERE id = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comment, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByName[models.Comment],
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCommentNotFound
		}

		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) GetByRecipeID(ctx context.Context, recipeID uuid.UUID) ([]models.Comment, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			id,
			recipe_id,
			author_id,
			body,
			rating,
			created_at,
			updated_at
		 FROM comments
		 WHERE recipe_id = $1
		 ORDER BY created_at DESC`,
		recipeID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[models.Comment],
	)

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO comments (
			recipe_id,
			author_id,
			body,
			rating
		 ) VALUES ($1,$2,$3,$4)
		 RETURNING id, created_at, updated_at`,
		comment.RecipeID,
		comment.AuthorID,
		comment.Body,
		comment.Rating,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if err := updateRecipeStats(ctx, comment.RecipeID, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *CommentRepository) Update(ctx context.Context, comment *models.Comment) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		`UPDATE comments
		 SET
			body = $2,
			rating = $3,
			updated_at = NOW()
		 WHERE id = $1
		 RETURNING recipe_id, updated_at`,
		comment.ID,
		comment.Body,
		comment.Rating,
	).Scan(&comment.RecipeID, &comment.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCommentNotFound
		}

		return err
	}

	if err := updateRecipeStats(ctx, comment.RecipeID, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *CommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var recipeID uuid.UUID
	err = tx.QueryRow(ctx,
		`DELETE FROM comments
		 WHERE id = $1
		 RETURNING recipe_id`,
		id,
	).Scan(&recipeID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCommentNotFound
		}
		return err
	}

	if err := updateRecipeStats(ctx, recipeID, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
