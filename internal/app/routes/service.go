package routes

import (
	"net/http"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/controllers"
)

func RegisterServiceRoutes(r *Router, ctl *controllers.HealthController) {
	serviceRoutes := []Route{
		{
			Path: "/healthcheck",
			Methods: map[string]Handler{
				http.MethodGet: {
					Func: ctl.HealthCheck,
				},
			},
		},
	}

	r.routes = append(r.routes, serviceRoutes...)
}
