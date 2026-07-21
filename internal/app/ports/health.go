package ports

import "context"

type HealthProbe interface {
	PingDB(ctx context.Context) error
}
