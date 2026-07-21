package routes

import (
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers"
)

func RegisterServiceRoutes(r *Router, ctl *controllers.HealthController) {
	serviceRoutes := []Route{
		{
			Path: "/healthcheck",
			Methods: map[string]http.HandlerFunc{
				http.MethodGet: ctl.HealthCheck,
			},
		},
	}

	r.routes = append(r.routes, serviceRoutes...)
}
