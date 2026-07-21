package services

import (
	"context"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/ports"
)

type HealthCheckService struct {
	probe ports.HealthProbe
}

func NewHealthCheckService(probe ports.HealthProbe) *HealthCheckService {
	return &HealthCheckService{probe: probe}
}

func (hcs *HealthCheckService) HealthCheck(ctx context.Context) (statusArray map[string]bool, ok bool) {
	statusArray = map[string]bool{
		"db": true,
	}
	ok = true

	if err := hcs.probe.PingDB(ctx); err != nil {
		statusArray["db"] = false
		ok = false
	}

	return statusArray, ok
}
