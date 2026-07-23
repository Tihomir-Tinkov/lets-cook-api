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

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Type() reflect.Type {
	return reflect.TypeOf(models.User{})
}

func (r *UserRepository) DB() *pgxpool.Pool {
	return r.db
}

func (r *UserRepository) Model() models.Model {
	return models.User{}
}

func (r *UserRepository) Find(
	ctx context.Context,
	id uuid.UUID,
) (models.User, error) {

	query := `
		SELECT
			id,
			display_name,
			email,
			password_hash,
			role,
			created_at,
			updated_at
		FROM users
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	row := r.db.QueryRow(ctx, query, args)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.DisplayName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
