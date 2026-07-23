package repositories

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *ImageRepository) Create(ctx context.Context, i *models.Image) error {
	var pgID pgtype.UUID
	copy(pgID.Bytes[:], i.ID[:])
	pgID.Valid = true
	_, err := r.db.Exec(ctx,
		`INSERT INTO images (id, filename, mime_type, extension, size, created_at)
		 VALUES($1,$2,$3,$4,$5,$6)`,
		pgID, i.FileName, i.MimeType, i.Extension, i.Size, time.Now(),
	)
	return err
}

func (r *ImageRepository) Get(ctx context.Context, id uuid.UUID) (*models.Image, error) {
	var i models.Image
	var pgID pgtype.UUID
	err := r.db.QueryRow(ctx,
		`SELECT id, filename, mime_type, extension, size FROM images WHERE id=$1`, id).
		Scan(&pgID, &i.FileName, &i.MimeType, &i.Extension, &i.Size)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("image not found")
		}
		return nil, err
	}
	i.ID = pgID.Bytes
	return &i, nil
}

func (r *ImageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var pgID pgtype.UUID
	copy(pgID.Bytes[:], id[:])
	pgID.Valid = true
	_, err := r.db.Exec(ctx, `DELETE FROM images WHERE id=$1`, pgID)
	return err
}
