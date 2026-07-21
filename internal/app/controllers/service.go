package controllers

import (
	"errors"
	"net/http"
	"os"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/responses"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/services"
)

type HealthController struct {
	healthCheckService *services.HealthCheckService
}

func NewServiceController(healthCheckService *services.HealthCheckService) *HealthController {
	return &HealthController{
		healthCheckService,
	}
}

// @Summary Health check
// @Description Health check
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Failure 500 {object} string
// @Router /healthcheck [get]
func (s *HealthController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	statusArray, ok := s.healthCheckService.HealthCheck(r.Context())
	if !ok {
		responses.JSONError(w, r, errors.New("some services failed"), http.StatusInternalServerError)
		return
	}

	responses.JSONResponse(w, http.StatusOK, map[string]any{
		"service": "healthcheck",
		"version": os.Getenv("VERSION"),
		"status":  statusArray,
	})
}
