package repositories

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseRepository interface {
	Type() reflect.Type
	DB() *pgxpool.Pool
	Model() models.Model
}

type RepositoryInterface[T models.Model] interface {
	Find(ctx context.Context, id uuid.UUID) (T, error)
	GetAllByID(ctx context.Context, ids []uuid.UUID) ([]T, error)
	FindAll(ctx context.Context, conditions []ports.SimpleCondition) ([]T, error)
	Model() T
}

type SimpleCondition = ports.SimpleCondition

func buildWhereClauseSimple(conditions []ports.SimpleCondition) (string, pgx.NamedArgs) {
	var whereClause []string
	args := pgx.NamedArgs{}

	for _, v := range conditions {
		whereClause = append(whereClause, v.Field+" = @"+v.Field)
		args[v.Field] = v.Value
	}

	return strings.Join(whereClause, " AND "), args
}

func runQuery(ctx context.Context, condition []ports.SimpleCondition, r BaseRepository, query string, orderBy string, limit int) (pgx.Rows, error) {
	args := pgx.NamedArgs{}
	var whereQuery string

	if condition != nil {
		whereQuery, args = buildWhereClauseSimple(condition)
		query += " WHERE " + whereQuery
	}

	if orderBy != "" {
		query += " ORDER BY " + orderBy
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.DB().Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
