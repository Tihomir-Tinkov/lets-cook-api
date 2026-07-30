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

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) DB() *pgxpool.Pool {
	return r.db
}

func (r *UserRepository) Model() *models.User {
	return &models.User{}
}

func (r *UserRepository) Type() reflect.Type {
	return reflect.TypeOf(models.User{})
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (
			display_name,
			email,
			password_hash,
			role
		 ) VALUES ($1,$2,$3,$4)
		 RETURNING id`,
		user.DisplayName,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT *
		 FROM users
		 WHERE id = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByName[models.User],
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT *
		 FROM users
		 WHERE email = $1`,
		email,
	)

	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByName[models.User],
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	result, err := r.db.Exec(
		ctx,
		`UPDATE users
		 SET
			display_name = $2,
			email = $3,
			password_hash = $4,
			updated_at = NOW()
		 WHERE id = $1`,
		user.ID,
		user.DisplayName,
		user.Email,
		user.PasswordHash,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM users
		 WHERE id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
