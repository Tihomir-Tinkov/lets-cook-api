package gateways

import (
	"context"
	"errors"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Health struct {
	dbPG *pgxpool.Pool
}

func NewHealth(dbPG *pgxpool.Pool) *Health {
	return &Health{
		dbPG: dbPG,
	}
}

func (h *Health) PingDB(ctx context.Context) error {
	if h.dbPG == nil {
		return errUnavailable
	}
	return h.dbPG.Ping(ctx)
}

var errUnavailable = errors.New("health probe unavailable")

var _ ports.HealthProbe = (*Health)(nil)
