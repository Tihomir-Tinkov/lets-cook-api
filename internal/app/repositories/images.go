package repositories

import (
	"context"
	"errors"
	"reflect"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrImageNotFound = errors.New("image not found")

type ImageRepository struct {
	db *pgxpool.Pool
}

func NewImageRepository(db *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) DB() *pgxpool.Pool {
	return r.db
}

func (r *ImageRepository) Model() *models.Image {
	return &models.Image{}
}

func (r *ImageRepository) Type() reflect.Type {
	return reflect.TypeOf(models.Image{})
}

func (r *ImageRepository) Create(ctx context.Context, image *models.Image) error {
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO images (
			filename,
			mime_type,
			extension,
			size
		) VALUES ($1,$2,$3,$4)
		RETURNING id, created_at`,
		image.FileName,
		image.MimeType,
		image.Extension,
		image.Size,
	).Scan(
		&image.ID,
		&image.CreatedAt,
	)

	return err
}

func (r *ImageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Image, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			id,
			filename,
			mime_type,
			extension,
			size,
			created_at
		 FROM images
		 WHERE id = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	image, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByName[models.Image],
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrImageNotFound
		}

		return nil, err
	}

	return &image, nil
}

func (r *ImageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM images
		 WHERE id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrImageNotFound
	}

	return nil
}
